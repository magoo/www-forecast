package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"fmt"
)

type Auth struct {
	*revel.Controller
}

func (c Auth) Create(googleIdToken string) revel.Result {
	ti, hd, err := models.VerifyIdToken(googleIdToken)

	if (err != nil){
		c.Flash.Error("Could not verify token from Google.")
		fmt.Println("Can't verify hd claim from Google: " + err.Error())
		return c.Render()
	}

	c.Session["user"] = ti.Email

	// Set the hosted domain. This is eventually a core privacy barrier.
	// If not a gsuite customer, it's public.
	// If gsuite customer, it's for the "hosted domain" only.
	// Cheap and simple privacy for the time being.
	// We change empty strings to "public" because dynamo doesn't like empty
	// strings downstream.
	if (hd == "") {
		c.Session["hd"] = "public"
	} else {
		c.Session["hd"] = hd
	}

	res := JSONResponse{Code: "ok"}
	return c.RenderJSON(res)
}

func (c Auth) Delete() revel.Result {
	c.Session["user"] = ""
	c.Session["hd"]		= ""
	c.Flash.Success("Logged Out.")

	res := JSONResponse{Code: "ok"}
	return c.RenderJSON(res)
}
