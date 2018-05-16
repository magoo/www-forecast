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

func (c Estimate) Update(eid string, title string, description string, unit string) revel.Result {

		models.UpdateEstimate(eid , title, description, unit, c.Session["user"])

		//Show success and redirect to the estimate w/ changes
		c.Flash.Success("Updated.")

		return c.Redirect("/view/estimate/%s", eid)

}
