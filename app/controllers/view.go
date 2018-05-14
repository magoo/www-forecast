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

			return c.Redirect(List.Index)
		}

		f := models.ViewScenario(sid)
		u :=  c.Session["user"]
		myForecast := models.ViewUserScenarioResults(u, sid)

		return c.Render(f, u, myForecast)
}

func (c View) Conclude(sid string, resultIndex int) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot conclude. Invalid scenario ID.")

		return c.Redirect(List.Index)

	}

	t := time.Now()

	f := models.ViewScenario(sid)

	if f.Owner != c.Session["user"] {
		c.Flash.Error("Cannot conclude scenario you do not own.")
		return c.Redirect(List.Index)
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

	f.Concluded = true
	f.ConcludedTime = t.String()
	f.Results = af
	f.ResultIndex = resultIndex
	f.BrierScore = bs

	models.PutItem(f, "scenarios-tf")

	return c.Redirect("/view/%s", sid)
}
