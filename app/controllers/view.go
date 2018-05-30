package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	"time"
	"fmt"
)

type View struct {
	*revel.Controller
}

func (c View) Index(sid string) revel.Result {

		c.Validation.Required(sid)
		c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Cannot view. Invalid scenario ID.")

			return c.Redirect(Home.List)
		}

		f := models.ViewScenario(sid)
		u :=  c.Session["user"]
		myForecast := models.ViewUserScenarioResults(u, sid)

		return c.Render(f, u, myForecast)
}

func (c View) Update(sid string, title string, description string, options []string) revel.Result{

	models.UpdateScenario(sid, title, description, options, c.Session["user"])
	c.Flash.Success("Updated.")

	return c.Redirect("/view/%s", sid)

}

func (c View) Conclude(sid string, resultIndex int) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot conclude. Invalid scenario ID.")

		return c.Redirect(Home.List)

	}

	t := time.Now()

	s := models.ViewScenario(sid)

	if s.Question.OwnerID != c.Session["user"] {
		c.Flash.Error("Cannot conclude scenario you do not own.")
		return c.Redirect(Home.List)
	}

	sr := models.ViewScenarioResults(sid)

	if (len(sr)== 0) {
		c.Flash.Error("No results to conclude!")
		return c.Redirect("/view/%s", sid)
	}

	fmt.Println("concluding results:", len(sr))
	af := getAverageForecasts(sr)

	//Calculate Brier Score
	bs := models.BrierCalc(af, resultIndex)

	s.Question.Concluded = true
	s.Question.ConcludedTime = t.String()
	s.Results = af
	s.ResultIndex = resultIndex
	s.Question.BrierScore = bs

	models.PutItem(s, "scenarios-tf")

	return c.Redirect("/view/%s", sid)
}

func (c View) AddRecord(sid string) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot view. Invalid scenario ID.")

		return c.Redirect(Home.List)
	}

	s := models.ViewScenario(sid)
	u :=  c.Session["user"]
	s.AddRecord(u)

	return c.Redirect("/view/%s", sid)


}
