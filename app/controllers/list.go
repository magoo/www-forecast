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
	f := models.ListScenarios(c.Session["hd"])
	if (len(f) == 0) {

		c.Flash.Error("There are no scenarios created yet.")

	}
	return c.Render(f)
}
