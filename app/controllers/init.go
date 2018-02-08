package controllers

import (
	"www-forecast/app/models"
	"github.com/revel/revel"
)
func init() {

	//Main auth. In all controllers, make sure the user is logged in.
	revel.InterceptFunc(checkUser, revel.BEFORE, &View{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &App{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Create{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Results{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Forecast{})

	revel.OnAppStart(models.DbConnect)

}

// Check for session token
func checkUser(c *revel.Controller) revel.Result {
		if (c.Session["user"] == ""){
			c.Flash.Error("Please log in first")
			return c.Redirect(Auth.Index)
		}
    return nil
}
