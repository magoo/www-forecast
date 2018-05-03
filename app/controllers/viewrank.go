package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	//"time"
)

type ViewRank struct {
	*revel.Controller
}

func (c ViewRank) Index(rid string) revel.Result {

		c.Validation.Required(rid)
		c.Validation.Match(rid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Cannot view. Invalid estimate ID.")

			return c.Redirect(List.Index)
		}

		r := models.GetRank(rid)
		u :=  c.Session["user"]

		return c.Render(r, u)
}

func (c ViewRank) Conclude(rid string, resultValue float64) revel.Result {

	c.Validation.Required(rid)
	c.Validation.Match(rid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
	//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$



	return c.Redirect("/view/estimate/%s", rid)
}
