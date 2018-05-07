package controllers

import (
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
	revel.InterceptFunc(checkUser, revel.BEFORE, &Estimate{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &ViewEstimate{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &EstimateResults{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Range{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Rank{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &Sort{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &ViewRank{})
	revel.InterceptFunc(checkUser, revel.BEFORE, &RankResults{})


	//revel.OnAppStart(models.DbConnect)

}

// Check for session token
func checkUser(c *revel.Controller) revel.Result {

	c.Validation.Required(c.Session["user"])

	if c.Validation.HasErrors() {

		//Redirect from unauthenticated link.
		c.Session["redirect"] = c.Request.URL.Path
		c.Flash.Error("Please login. You'll be redirected to the URL you were trying to visit.")

		return c.Redirect(Auth.Index)
	}

    return nil
}
