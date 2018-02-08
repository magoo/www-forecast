package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type Create struct {
	*revel.Controller
}

func (c Create) Index() revel.Result {
		return c.Render()
}

func (c Create) Create(title string, description string, options []string) revel.Result {
		sid := models.CreateScenario(title, description, options, c.Session["hd"])
		//fmt.Println(options[0])
		return c.Redirect("/view/%s", sid )
}
