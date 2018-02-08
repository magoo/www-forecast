package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	f := models.ListScenarios(c.Session["hd"])
	return c.Render(f)
}
