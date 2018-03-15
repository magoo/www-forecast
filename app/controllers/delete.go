package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
)

type Delete struct {
	*revel.Controller
}


func (c Delete) DeleteScenario(sid string) revel.Result {
		c.Validation.Match(sid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))

		if c.Validation.HasErrors() {
			c.Flash.Error("You have to identify a scenario.")
			return c.Redirect(Create.Index)
		}


		models.DeleteScenario(sid, c.Session["user"])

		res := JSONResponse{Code: "ok"}

		return c.RenderJSON(res)
}
