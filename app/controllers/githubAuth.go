package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"

	"github.com/google/go-github/github"

	"strconv"

	"os"

	"golang.org/x/oauth2"
)

type GithubAuth struct {
	*revel.Controller
}

// TODO: Dry out with init.go
type GithubOauthConfig struct {
	GithubClientID     string
	GithubClientSecret string
}

var (
	GithubConfig = &GithubOauthConfig{
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}

	GithubEndpoint = oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize?scope=user:email",
		TokenURL: "https://github.com/login/oauth/access_token",
	}
)

func (c GithubAuth) Callback() revel.Result {
	GithubOauth2Config := &oauth2.Config{
		ClientID:     GithubConfig.GithubClientID,
		ClientSecret: GithubConfig.GithubClientSecret,
		RedirectURL:  revel.Config.StringDefault("e6eDomain", "https://localhost:9000") + "/github/callback",
		Endpoint:     GithubEndpoint,
	}

	ctx := c.Request.Context()

	req := c.Request
	// TODO: Include MIT license with this code taken from the library
	var cookieConfig gologin.CookieConfig
	if revel.Config.BoolDefault("mode.dev", false) {
		cookieConfig = gologin.DebugOnlyCookieConfig
	} else {
		cookieConfig = gologin.DefaultCookieConfig
	}

	//// Get the OAuth state value from the cookie
	ownerState := GetOAuthState(cookieConfig, req, ctx)

	//// Compare the state value and get the Auth token from the authorization code
	token := GetToken(GithubOauth2Config, req, ctx, ownerState)

	githubUser := GetGithubUser(GithubOauth2Config, req, ctx, token)

	//// Get primary email address for github account from API
	emailReq, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)

	httpClient := GithubOauth2Config.Client(ctx, token)
	resp, err := httpClient.Do(emailReq)

	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	type UserEmails []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}
	var userEmails UserEmails
	_ = json.Unmarshal([]byte(body), &userEmails)

	var primaryEmail string
	for _, n := range userEmails {
		if n.Primary {
			primaryEmail = n.Email
		}
	}

	//// See if this user already exists in the database
	oAuthId := strconv.FormatInt(*githubUser.ID, 10)
	user, needs_creating := models.GetUserByOAuth(oAuthId, "github")
	if needs_creating {
		// Save the email address and provider in the DB
		models.CreateUser(primaryEmail, oAuthId, "github")
		user, _ = models.GetUserByOAuth(oAuthId, "github")
		fmt.Println("Saving user to DB: " + user.Email)
	} else {
		fmt.Println("Found user in DB: " + user.Email)
	}

	// Set the user ID on the user session
	c.Session["user"] = "github:" + user.OauthID
	return c.Redirect(Home.List)
}

// Mimics gologin's StateHandler
func GetOAuthState(config gologin.CookieConfig, req *revel.Request, ctx context.Context) string {
	cookie, err := req.Cookie(config.Name)
	if err == nil {
		return cookie.GetValue()
	} else {
		return ""
	}
}

// Mimics oauth2.parseCallback
func parseCallback(req *revel.Request) (authCode, state string, err error) {
	err = req.ParseForm()
	if err != nil {
		return "", "", err
	}
	authCode = req.Form.Get("code")
	state = req.Form.Get("state")
	if authCode == "" || state == "" {
		fmt.Println("authCode", authCode) // revertme
		fmt.Println("state", state)       // revertme
		return "", "", errors.New("oauth2: Request missing code or state")
	}
	return authCode, state, nil
}

// Mimics gologin's oauth2.CallbackHandler
func GetToken(config *oauth2.Config, req *revel.Request, ctx context.Context, ownerState string) *oauth2.Token {
	authCode, state, err := parseCallback(req)
	if err != nil {
		fmt.Println("Error getting authCode from callback request:", err)
		return nil
	}
	if state != ownerState || state == "" {
		fmt.Println("oauth2 state parameter:", state) // revertme
		fmt.Println("ownerState:", ownerState)        // revertme
		fmt.Println("State and ownerstate don't match:", state, ownerState)
		return nil
	}
	// use the authorization code to get a Token
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		fmt.Println("Error getting token from auth code:", err)
		return nil
	}

	return token
}

func GetGithubUser(config *oauth2.Config, req *revel.Request, ctx context.Context, token *oauth2.Token) *github.User {
	httpClient := config.Client(ctx, token)

	githubClient := github.NewClient(httpClient)

	user, resp, err := githubClient.Users.Get(ctx, "")
	err = validateResponse(user, resp, err)
	if err != nil {
		fmt.Println("Error validating github response:", err)
		return nil
	}

	return user
}

// validateResponse returns an error if the given Github user, raw
// http.Response, or error are unexpected. Returns nil if they are valid.
func validateResponse(user *github.User, resp *github.Response, err error) error {
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.New("github: unable to get Github User")
	}
	if user == nil || user.ID == nil {
		return errors.New("github: unable to get Github User")
	}
	return nil
}
