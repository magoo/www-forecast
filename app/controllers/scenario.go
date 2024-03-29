package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Scenario struct {
	*revel.Controller
}

type JSONResponse struct {
	Code string `json:"code"`
}

func (c Scenario) Index() revel.Result {
	return c.Render()
}

func (c Scenario) Create(title string, description string, options []string) revel.Result {

	c.Validation.MinSize(title, 1)

	for _, o := range options {
		c.Validation.MinSize(o, 1)
	}

	if c.Validation.HasErrors() {
		c.Flash.Error("You can't have an empty title or option.")
		return c.Redirect(Scenario.Create)
	}

	sid := models.CreateScenario(title, description, options, c.Session["hd"].(string), c.Session["user"].(string))
	//fmt.Println(options[0])

	c.Flash.Out["createdurl"] = revel.Config.StringDefault("e6eDomain", "https://www.e6e.io") + "/view/scenario/" + sid

	return c.Redirect(Home.List)

}

func (c Scenario) Delete(id string) revel.Result {
	c.Validation.Match(id, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))

	if c.Validation.HasErrors() {
		c.Flash.Error("You have to identify a scenario.")
		return c.Redirect(Scenario.Create)
	}

	models.DeleteScenario(id, c.Session["user"].(string))

	res := JSONResponse{Code: "ok"}

	return c.RenderJSON(res)
}

func (c Scenario) View(sid string) revel.Result {

	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view. Invalid scenario ID.")

		return c.Redirect(Home.List)
	}

	f := models.ViewScenario(sid)
	u := c.Session["user"].(string)
	myForecast := models.ViewUserScenarioResults(u, sid)

	return c.Render(f, u, myForecast)
}

func (c Scenario) Update(sid string, title string, description string, options []string) revel.Result {

	models.UpdateScenario(sid, title, description, options, c.Session["user"].(string))
	c.Flash.Success("Updated.")

	return c.Redirect("/view/scenario/%s", sid)

}

func (c Scenario) Conclude(sid string, resultIndex int) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot conclude. Invalid scenario ID.")

		return c.Redirect(Home.List)

	}

	t := time.Now()

	s := models.ViewScenario(sid)

	if s.Question.OwnerID != c.Session["user"].(string) {
		c.Flash.Error("Cannot conclude scenario you do not own.")
		return c.Redirect(Home.List)
	}

	sr := models.ViewScenarioResults(sid)

	if len(sr) == 0 {
		c.Flash.Error("No results to conclude!")
		return c.Redirect("/view/scenario/%s", sid)
	}

	rd, _ := s.CalcResultData()

	//Calculate Brier Score
	bs := models.BrierCalc(rd.Avg, resultIndex)

	s.Question.Concluded = true
	s.Question.ConcludedTime = t.String()
	//s.Results = af
	s.ResultIndex = resultIndex

	//This is rolling brier score code. Removing.
	//if s.Question.BrierScore == 0 {
	//	s.Question.BrierScore = bs
	//} else {
	//	s.Question.BrierScore = (bs + s.Question.BrierScore) / 2
	//}

	s.Question.BrierScore = bs

	err := models.WriteQuestion(s)

	if err != nil {
		fmt.Println("Error writing question.")
	}

	u := c.Session["user"].(string)

	err = s.AddRecord(u)

	if err != nil {
		fmt.Println("Error writing record to scenario.")
	}

	err = s.ConcludeScenarioForecasts()

	if err != nil {
		fmt.Println("Error concluding individual forecasts.")
	}

	err = s.Question.WriteRecord("Concluded. Brier Score is updated to "+strconv.FormatFloat(s.Question.BrierScore, 'f', -1, 64), c.Session["user"].(string))

	if err != nil {
		fmt.Println("Error concluding scenario.")
	}

	// models.DeleteQuestionAnswers(sid)

	return c.Redirect("/view/scenario/%s", sid)
}

func (c Scenario) AddRecord(sid string) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view. Invalid scenario ID.")

		return c.Redirect(Home.List)
	}

	s := models.ViewScenario(sid)
	u := c.Session["user"].(string)
	err := s.AddRecord(u)

	if err != nil {
		c.Flash.Error("Nothing to record.")
	}

	c.Flash.Success("Results added to record.")
	return c.Redirect("/view/scenario/%s", sid)

}

func (c Scenario) Results(sid string) revel.Result {

	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view results, errors in submission.")

		return c.Redirect(Home.List)
	}

	//This attempts to retrieve the scenario. It needs to exist to continue.
	s := models.ViewScenario(sid)

	// We use the SID from the successful call using the hosted domain, instead of whatever the user gives us.
	sr := models.ViewScenarioResults(s.Question.Id)

	// Variable `sr` has an array of scenario results. 

	if len(sr) == 0 {
		c.Flash.Error("No results yet.")
		return c.Redirect("/view/scenario/%s", sid)
	}

	
	rd, err := s.CalcResultData()

	if err != nil {
		c.Flash.Error("Error calculating average forecast.")
	}

	

	return c.Render(sr, s, rd)
	 



}
