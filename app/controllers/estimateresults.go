package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	)

type EstimateResults struct {
	*revel.Controller
}

func (c EstimateResults) Index(eid string) revel.Result {

		c.Validation.Required(eid)
		c.Validation.Match(eid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Cannot view results, errors in submission.")

			return c.Redirect(List.Index)
		}

		//This attempts to retrieve the scenario based on the hosted domain, for security.
		e := models.GetEstimate(eid)

		// We use the SID from the successful call using the hosted domain, instead of whatever the user gives us.
		er := models.ViewEstimateResults(e.Eid)
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
