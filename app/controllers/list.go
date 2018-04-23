package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"

)

type List struct {
	*revel.Controller
}

func (c List) Index() revel.Result {

	if (c.Session["redirect"] != "") {

		redirect := c.Session["redirect"]
		delete(c.Session, "redirect") // Removed item from session
		return c.Redirect(redirect)
	}

	//The interceptor in init() should enforce that we have this.
	//This protects us just in case, enforcing literally anything in the "hd" field.
	//fmt.Println("App controller is launching")
	f := models.ListScenarios(c.Session["user"])

	empty := true

	if (len(f) > 0) {

		empty = false

	}

	e := models.ListEstimates(c.Session["user"])

	if (len(e) > 0) {

		empty = false

	}

	return c.Render(f, e, empty)
}
