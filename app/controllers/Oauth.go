package controllers

import (
	"github.com/revel/revel"
)

type Oauth struct {
	*revel.Controller
}

func Delete(c Oauth) revel.Result {
	c.Session["user"] = ""
	c.Session["hd"] = ""
	c.Flash.Success("Logged Out.")

	res := JSONResponse{Code: "ok"}
	return c.RenderJSON(res)
}
