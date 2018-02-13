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
		s := 0
		for _, v := range value {
			s += v
		}

		if s != 100 {
			c.Flash.Error("You need your values to equal 100%")
			return c.Redirect("/view/%s", sid )
		}

		models.CreateForecast(c.Session["user"], value, sid, c.Session["hd"])
		//fmt.Println(options[0])
		return c.Redirect("/view/%s/results", sid )
}
