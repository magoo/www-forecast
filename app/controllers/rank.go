package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type Rank struct {
	*revel.Controller
}

func (c Rank) Index() revel.Result {
		return c.Render()
}

func (c Rank) Create(title string, description string, options []string) revel.Result {

    rid := models.CreateRank(title, description, options, c.Session["hd"], c.Session["user"])

		return c.Redirect("/view/rank/%s", rid)
}

func (c Rank) Update(rid string, title string, description string, options []string) revel.Result {

	models.UpdateRank(rid, title, description, options, c.Session["user"])

	c.Flash.Success("Updated rank.")
	return c.Redirect("/view/rank/%s", rid)
}
