package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	)

type Results struct {
	*revel.Controller
}

func (c Results) Index(sid string) revel.Result {

		c.Validation.Required(sid)
		c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Invalid scenario ID.")

			return c.Redirect(List.Index)
		}

		//This attempts to retrieve the scenario based on the hosted domain, for security.
		s := models.ViewScenario(sid, c.Session["hd"])

		// We use the SID from the successful call using the hosted domain, instead of whatever the user gives us.
		sr := models.ViewScenarioResults(s.Sid, c.Session["hd"])
		if (len(sr)>0){
			avg := getAverageForecasts(sr)
			return c.Render(sr, s, avg)
		} else {
			c.Flash.Error("No results yet.")
			return c.Redirect("/view/%s", sid)
		}

}

func getAverageForecasts(sr []models.Forecast) ([]int){


	avg := []int{}
	size := len(sr[0].Forecasts)

	for i := 0; i < size; i++ {
		sum := 0
			for _, v := range sr {
					sum += v.Forecasts[i]
					//fmt.Println("Adding forecast: ", v.Forecasts[i])
			}
			//fmt.Println("Adding average to array: ", sum / len(sr))
		avg = append(avg, sum / len(sr))
	}

	return avg
}
