package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type Forecast struct {
	*revel.Controller
}

func (c Forecast) Index() revel.Result {
		return c.Render()
}

func (c Forecast) Create(value []int, sid string) revel.Result {

		models.CreateForecast(c.Session["user"], value, sid)
		//fmt.Println(options[0])
		return c.Redirect("/view/%s/results", sid )
}
