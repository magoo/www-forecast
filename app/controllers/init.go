package controllers

import (
	"github.com/revel/revel"
)
func init() {

	revel.InterceptFunc(enforceHSTS, revel.BEFORE, &Home{})

	//Main auth. In all controllers, make sure the user is logged in.
	//Every controller with sensitive content should be here.
	//Better yet, whitelisting these controllers would be better.
	revel.InterceptFunc(checkUser, revel.BEFORE, &Scenario{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Forecast{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Estimate{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Range{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Rank{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Sort{})

	//revel.OnAppStart(models.DbConnect)

	revel.TemplateFuncs["increment"] = func(a int) int {
		return a + 1
	}

}

func enforceHSTS(c *revel.Controller) revel.Result {

		if c.Request.Header.Get("X-Forwarded-Proto") != "" {
			if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
				return c.Redirect("https://e6e.io")
			}
		}

	return nil
}

// Check for session token
func checkUser(c *revel.Controller) revel.Result {

	revel.AppLog.Debug("AccessLog", "user", c.Session["user"], "ip", c.ClientIP, "path", c.Request.URL.Path)

	c.Validation.Required(c.Session["user"])

	if c.Validation.HasErrors() {

		//Redirect from unauthenticated link.
		c.Session["redirect"] = c.Request.URL.Path
		c.Flash.Error("Please login. You'll be redirected to the URL you were trying to visit.")

		return c.Redirect(Home.Index)
	}

    return nil
}
