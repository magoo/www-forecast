package controllers

import (
	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Forecast struct {
	*revel.Controller
}

func (c Forecast) Index() revel.Result {
	return c.Render()
}

func (c Forecast) Create(value []float64, sid string) revel.Result {
	s := float64(0)
	for _, v := range value {
		s += v
	}

	if s != 100 {
		c.Flash.Error("You need your values to equal 100%")
		return c.Redirect("/view/scenario/%s", sid)
	}

	models.CreateForecast(c.Session["user"].(string), value, sid, c.Session["hd"].(string))
	//fmt.Println(options[0])
	return c.Redirect("/view/scenario/%s/results", sid)
}
