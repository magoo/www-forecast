package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Estimate struct {
	*revel.Controller
}

func (c Estimate) Index() revel.Result {
	return c.Render()
}

func (c Estimate) Create(title string, description string, unit string) revel.Result {

	eid := models.CreateEstimate(title, description, unit, c.Session["hd"].(string), c.Session["user"].(string))

	c.Flash.Out["createdurl"] = revel.Config.StringDefault("e6eDomain", "https://www.e6e.io") + "/view/estimate/" + eid

	return c.Redirect(Home.List)
}

func (c Estimate) View(eid string) revel.Result {

	c.Validation.Required(eid)
	c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view. Invalid estimate ID.")

		return c.Redirect(Home.List)
	}

	e := models.GetEstimate(eid)
	u := c.Session["user"]

	return c.Render(e, u)
}

func (c Estimate) Record(eid string) revel.Result {
	c.Validation.Required(eid)
	c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view. Invalid estimate ID.")

		return c.Redirect(Home.List)
	}

	e := models.GetEstimate(eid)
	u := c.Session["user"]
	err := e.AddRecord(u.(string))

	if err != nil {
		c.Flash.Error("Nothing to record.")
	} else {
		c.Flash.Success("Results added to record.")
	}
	return c.Redirect("/view/estimate/%s", eid)

}

func (c Estimate) Conclude(eid string, resultValue float64) revel.Result {

	c.Validation.Required(eid)
	c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot conclude. Invalid estimate ID.")

		return c.Redirect(Home.List)

	}

	e := models.GetEstimate(eid)

	if e.Question.OwnerID != c.Session["user"] {
		c.Flash.Error("Cannot conclude estimate you do not own.")
		return c.Redirect(Home.List)
	}

	er := models.ViewEstimateResults(eid)

	if len(er) == 0 {
		c.Flash.Error("No results to conclude!")
		return c.Redirect("/view/estimate/%s", eid)
	}

	emin, emax := models.GetAverageRange(er)

	//Calculate Brier Score
	//This is different. There is a 90% confidence assumption.
	bs := models.BrierCalcEstimate(emin, emax, resultValue)

	t := time.Now()
	e.Concluded = true
	e.ConcludedTime = t.String()
	//e.AvgMinimum = emin
	//e.AvgMaximum = emax
	//e.Actual = resultValue

	//if e.Question.BrierScore == 0 {
	//	e.Question.BrierScore = bs
	//
	//} else {
	//	e.Question.BrierScore = (bs + e.Question.BrierScore) / 2
	//}

	e.Question.BrierScore = bs

	//err := models.PutItem(e, questionTable)
	err := models.WriteQuestion(e)

	if err != nil {
		fmt.Println("Error writing question.")
	}

	u := c.Session["user"]
	err = e.AddRecord(u.(string))

	if err != nil {
		fmt.Println("Error writing record to question.")
		return c.Redirect("/view/estimate/%s", eid)
	}

	err = e.Question.WriteRecord("Concluded. Brier Score is updated to "+strconv.FormatFloat(e.Question.BrierScore, 'f', -1, 64), c.Session["user"].(string))

	// models.DeleteQuestionAnswers(eid)

	err = e.ConcludeEstimateResults()

	if err != nil {
		fmt.Println("Error concluding individual ranges.")
	}

	if err != nil {
		fmt.Println("Error concluding question.")
	}
	c.Flash.Success("Updated score.")
	return c.Redirect("/view/estimate/%s", eid)
}

func (c Estimate) Update(eid string, title string, description string, unit string) revel.Result {

	models.UpdateEstimate(eid, title, description, unit, c.Session["user"].(string))

	//Show success and redirect to the estimate w/ changes
	c.Flash.Success("Updated.")

	return c.Redirect("/view/estimate/%s", eid)

}

func (c Estimate) Delete(id string) revel.Result {
	c.Validation.Match(id, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))

	if c.Validation.HasErrors() {
		c.Flash.Error("You have to identify an estimate.")
		return c.Redirect(Home.List)
	}

	models.DeleteEstimate(id, c.Session["user"].(string))

	res := JSONResponse{Code: "ok"}

	return c.RenderJSON(res)
}

func (c Estimate) Results(eid string) revel.Result {

	c.Validation.Required(eid)
	c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view results, errors in submission.")

		return c.Redirect(Home.List)
	}

	//This attempts to retrieve the scenario based on the hosted domain, for security.
	e := models.GetEstimate(eid)

	// We use the SID from the successful call using the hosted domain, instead of whatever the user gives us.
	er := models.ViewEstimateResults(e.Question.Id)
	if len(er) > 0 {
		avgmin, avgmax := models.GetAverageRange(er)
		return c.Render(er, e, avgmin, avgmax)
	} else {
		c.Flash.Error("No results yet.")
		return c.Redirect("/view/estimate/%s", eid)
	}

}
