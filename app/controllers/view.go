package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
)

type View struct {
	*revel.Controller
}

func (c View) Index(sid string) revel.Result {

		c.Validation.Required(sid)
		c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Invalid scenario ID.")

			return c.Redirect(List.Index)
		}

		f := models.ViewScenario(sid, c.Session["hd"])
		return c.Render(f)
}
