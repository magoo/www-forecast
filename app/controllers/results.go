package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"fmt"
)

type Results struct {
	*revel.Controller
}

func (c Results) Index(fid string) revel.Result {
		f := models.ViewForecast(fid, c.Session["hd"])
		fr := models.ViewForecastResults(fid)
		avg := getAverageForecasts(fr)
		return c.Render(fr, f, avg)
}

func getAverageForecasts(fr []models.Cast) ([]int){

	avg := []int{}
	size := len(fr[0].Forecasts)

	for i := 0; i < size; i++ {
		sum := 0
			for _, v := range fr {
					sum += v.Forecasts[i]
					fmt.Println("Adding forecast: ", v.Forecasts[i])
			}
			fmt.Println("Adding average to array: ", sum / len(fr))
		avg = append(avg, sum / len(fr))
	}

	return avg
}
