package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	"time"
)

type Estimate struct {
	*revel.Controller
}

func (c Estimate) Index() revel.Result {
		return c.Render()
}

func (c Estimate) Create(title string, description string, unit string) revel.Result {

    eid := models.CreateEstimate(title, description, unit, c.Session["hd"], c.Session["user"])

		return c.Redirect("/view/estimate/%s", eid)
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
		u :=  c.Session["user"]

		return c.Render(e, u)
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

	models.PutItem(e, "questions-tf")

	return c.Redirect("/view/estimate/%s", eid)
}

func (c Estimate) Update(eid string, title string, description string, unit string) revel.Result {

		models.UpdateEstimate(eid , title, description, unit, c.Session["user"])

		//Show success and redirect to the estimate w/ changes
		c.Flash.Success("Updated.")

		return c.Redirect("/view/estimate/%s", eid)

}

func (c Estimate) Delete(eid string) revel.Result {
		c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))

		if c.Validation.HasErrors() {
			c.Flash.Error("You have to identify an estimate.")
			return c.Redirect(Home.List)
		}


		models.DeleteEstimate(eid, c.Session["user"])

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
		if (len(er)>0){
			avgmin, avgmax := getAverageRange(er)
			return c.Render(er, e, avgmin, avgmax)
		} else {
			c.Flash.Error("No results yet.")
			return c.Redirect("/view/estimate/%s", eid)
		}

}

func getAverageRange(er []models.Range) (avgmin float64, avgmax float64){

	size := len(er)
	var sum float64 = 0

	for _, v := range er {
		sum += v.Minimum
	}

	avgmin = sum / float64(size)

	//reset
	sum = 0

	for _, v := range er {
		sum += v.Maximum
	}

	avgmax = sum / float64(size)

	return
}
