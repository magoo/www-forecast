package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	)

type Results struct {
	*revel.Controller
}

func (c Results) Index(sid string) revel.Result {
		s := models.ViewScenario(sid, c.Session["hd"])
		sr := models.ViewScenarioResults(sid)
		avg := getAverageForecasts(sr)
		return c.Render(sr, s, avg)
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
