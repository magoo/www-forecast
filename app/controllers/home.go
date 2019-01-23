package controllers

import (
	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Home struct {
	*revel.Controller
}

func (c Home) Index() revel.Result {
	return c.Render()
}

func (c Home) List() revel.Result {

	if c.Session["user"] == "" {
		c.Flash.Error("Please log in.")
		return c.Redirect(Home.Index)
	}

	if c.Session["redirect"] != "" {

		redirect := c.Session["redirect"]
		delete(c.Session, "redirect") // Removed item from session
		return c.Redirect(redirect)
	}

	//The interceptor in init() should enforce that we have this.
	//This protects us just in case, enforcing literally anything in the "hd" field.
	//fmt.Println("App controller is launching")
	qs := models.ListQuestions(c.Session["user"])

	empty := true

	if len(qs) > 0 {

		empty = false

	}

	//Eventually show recent answers
	//	as := models.ListAnswers(c.Session["user"])
	//	if (len(as) > 0) {
	//		empty = false
	//	}

	return c.Render(qs, empty)
}

func (c Home) Alpha() revel.Result {
	return c.Render()
}
