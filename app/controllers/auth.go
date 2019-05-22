package controllers

import (
	"fmt"

	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Auth struct {
	*revel.Controller
}

func (c Auth) Create(googleIdToken string) revel.Result {
	ti, hd, err := models.VerifyIdToken(googleIdToken)

	if err != nil {
		c.Flash.Error("Could not verify token from Google.")
		fmt.Println("Can't verify hd claim from Google: " + err.Error())
		return c.Render()
	}

	user, needs_creating := models.GetUserByOAuth(ti.UserId, "google")
	if needs_creating {
		models.SaveUser(ti.Email, ti.UserId, "google")
		user, _ = models.GetUserByOAuth(ti.UserId, "google")
	}

	c.Session["user"] = user.Id

	// Set the hosted domain. This is eventually a core privacy barrier.
	// If not a gsuite customer, it's public.
	// If gsuite customer, it's for the "hosted domain" only.
	// Cheap and simple privacy for the time being.
	// We change empty strings to "public" because dynamo doesn't like empty
	// strings downstream.
	if hd == "" {
		c.Session["hd"] = "public"
	} else {
		c.Session["hd"] = hd
	}

	res := JSONResponse{Code: "ok"}
	return c.RenderJSON(res)
}

func (c Auth) Delete() revel.Result {
	c.Session["user"] = ""
	c.Session["hd"] = ""
	c.Flash.Success("Logged Out.")

	res := JSONResponse{Code: "ok"}
	return c.RenderJSON(res)
}
