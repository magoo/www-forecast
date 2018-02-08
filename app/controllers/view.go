package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type View struct {
	*revel.Controller
}

func (c View) Index(sid string) revel.Result {
		f := models.ViewScenario(sid, c.Session["hd"])
		return c.Render(f)
}
