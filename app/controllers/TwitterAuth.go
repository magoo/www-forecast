package controllers

import (
	"fmt"
	"os"

	// "github.com/dghubble/gologin/twitter"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"
)

type TwitterAuth struct {
	*revel.Controller
}

func (c TwitterAuth) Callback() revel.Result {
	oauthToken := c.Params.Query.Get("oauth_token")
	verifier := c.Params.Query.Get("oauth_verifier")

	twitterConfig := &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CLIENT_ID"),
		ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		CallbackURL:    "http://localhost:9000/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	accessToken, accessSecret, err := twitterConfig.AccessToken(oauthToken, "", verifier)
	if err != nil {
		return c.RenderError(err)
	}
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := twitterConfig.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)

	accountVerifyParams := &twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(true),
	}
	twitterUser, _, err := twitterClient.Accounts.VerifyCredentials(accountVerifyParams)
        if err != nil {
                return c.RenderError(err)
        }

	oAuthId := twitterUser.IDStr
	email := twitterUser.Email
	user, needs_creating := models.GetUserByOAuth(oAuthId, "twitter")
	if needs_creating {
		// Save the email address and provider in the DB
		models.CreateUser(email, oAuthId, "twitter")
		user, _ = models.GetUserByOAuth(oAuthId, "twitter")
		fmt.Println("Saving user to DB: " + email)
	} else {
		fmt.Println("Found user in DB: " + email)
	}

	c.Session["user"] = "twitter:" + user.OauthID
	return c.Redirect(Home.List)

}
