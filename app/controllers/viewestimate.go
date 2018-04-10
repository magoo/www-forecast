package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	"time"
)

type ViewEstimate struct {
	*revel.Controller
}

func (c ViewEstimate) Index(eid string) revel.Result {

		c.Validation.Required(eid)
		c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Cannot view. Invalid estimate ID.")

			return c.Redirect(List.Index)
		}

		e := models.GetEstimate(eid)
		u :=  c.Session["user"]

		return c.Render(e, u)
}

func (c ViewEstimate) Conclude(eid string, resultValue float64) revel.Result {

	c.Validation.Required(eid)
	c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

	if c.Validation.HasErrors() {
		c.Flash.Error("Cannot conclude. Invalid estimate ID.")

		return c.Redirect(List.Index)

	}

	e := models.GetEstimate(eid)

	if e.Owner != c.Session["user"] {
		c.Flash.Error("Cannot conclude estimate you do not own.")
		return c.Redirect(List.Index)
	}

	er := models.ViewEstimateResults(eid)
	t := time.Now()

	if (len(er)== 0) {
		c.Flash.Error("No results to conclude!")
		return c.Redirect("/view/estimate/%s", eid)
	}

	emin, emax := getAverageRange(er)

	//Calculate Brier Score
	//This is different. There is a 90% confidence assumption.
	bs := models.BrierCalcEstimate(emin, emax, resultValue)

	e.Concluded = true
	e.ConcludedTime = t.String()
	e.AvgMinimum = emin
	e.AvgMaximum = emax
	e.Actual = resultValue
	e.BrierScore = bs

	models.PutItem(e, "estimates-tf")

	return c.Redirect("/view/estimate/%s", eid)
}
