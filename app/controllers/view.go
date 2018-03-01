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
			c.Flash.Error("Invalid scenario ID.")

			return c.Redirect(List.Index)
		}

		f := models.ViewScenario(sid, c.Session["hd"])

		return c.Render(f)
}

func (c View) Conclude(sid string, resultIndex int) revel.Result {
	c.Validation.Required(sid)
	c.Validation.Required(resultIndex)
	c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Invalid scenario ID.")

		return c.Redirect(List.Index)

	}

	t := time.Now()

	f := models.ViewScenario(sid, c.Session["hd"])
	sr := models.ViewScenarioResults(sid, c.Session["hd"])

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

	models.PutItem(f, "scenarios")

	return c.Redirect("/view/%s", sid)
}
