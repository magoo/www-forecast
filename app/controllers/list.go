package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"

)

type List struct {
	*revel.Controller
}

func (c List) Index() revel.Result {

	//The interceptor in init() should enforce that we have this.
	//This protects us just in case, enforcing literally anything in the "hd" field.
	//fmt.Println("App controller is launching")
	f := models.ListScenarios(c.Session["user"])

	empty := false

	if (len(f) == 0) {

		empty = true

	}

	return c.Render(f, empty)
}
