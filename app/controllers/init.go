package controllers

import (
	"www-forecast/app/models"
	"github.com/revel/revel"
)
func init() {

	//Main auth. In all controllers, make sure the user is logged in.
	//Every controller with sensitive content should be here.
	//Better yet, whitelisting these controllers would be better.
	revel.InterceptFunc(checkUser, revel.BEFORE, &View{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &List{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Create{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Results{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Forecast{})

	revel.OnAppStart(models.DbConnect)

}

// Check for session token
func checkUser(c *revel.Controller) revel.Result {

	c.Validation.Required(c.Session["user"])

	if c.Validation.HasErrors() {
		c.Flash.Error("Please, login!")
		return c.Redirect(Auth.Index)
	}

    return nil
}
