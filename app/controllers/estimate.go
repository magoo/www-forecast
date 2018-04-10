package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type Estimate struct {
	*revel.Controller
}

func (c Estimate) Index() revel.Result {
		return c.Render()
}

func (c Estimate) Create(title string, description string, unit string) revel.Result {

    eid := models.CreateEstimate(title, description, unit, c.Session["hd"], c.Session["user"])

		return c.Redirect("/view/estimate/%s", eid)
}
