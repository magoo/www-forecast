package controllers

import (
	"github.com/revel/revel"

)

type Logout struct {
	*revel.Controller
}

func (c Logout) Index() revel.Result {

	c.Session["user"] = ""
	c.Session["hd"]		= ""
	c.Flash.Success("Logged Out.")

	return c.Redirect(Auth.Index)
}
