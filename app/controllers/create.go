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
		c.Validation.MinSize(title, 1)

		for _, o := range options{
			c.Validation.MinSize(o, 1)
		}


		if c.Validation.HasErrors() {
			c.Flash.Error("You can't have an empty title or option.")
			return c.Redirect(Create.Index)
		}


		sid := models.CreateScenario(title, description, options, c.Session["hd"])
		//fmt.Println(options[0])
		return c.Redirect("/view/%s", sid)
}
