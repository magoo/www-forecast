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

func (c Forecast) Create(value []int, fid string) revel.Result {

		models.CreateCast(c.Session["user"], value, fid)
		//fmt.Println(options[0])
		return c.Redirect("/view/%s", fid )
}
