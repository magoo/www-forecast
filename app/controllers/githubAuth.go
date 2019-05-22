package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"io/ioutil"
	"net/http"

	go_github "github.com/google/go-github/github"

	"strconv"

	"golang.org/x/oauth2"
	oauth2Login "github.com/dghubble/gologin/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
	"os"
)

//import (
//	"github.com/dghubble/gologin"
//	"github.com/dghubble/gologin/github"
//	"github.com/dghubble/sessions"
//	"golang.org/x/oauth2"
//	githubOAuth2 "golang.org/x/oauth2/github"
//)
//
//const (
//	sessionName    = "example-github-app" //nocommit
//	sessionSecret  = "example cookie signing secret" //nocommit
//	sessionUserKey = "githubID" //nocommit
//)


type GithubAuth struct {
	*revel.Controller
}

// TODO: Dry out with init.go
type GithubOauthConfig struct {
	GithubClientID     string
	GithubClientSecret string
}

func (c GithubAuth) Callback() revel.Result {
	fmt.Println("params", c.Params)
	fmt.Println("args", c.Args)
	fmt.Println("viewargs", c.ViewArgs)
	ctx := c.Request.Context()
	fmt.Println("context:", ctx)


	req := c.Request
	// TODO: Include MIT license with this code taken from the library
	// TODO: Dry out with init.go
	cookieConfig := gologin.DebugOnlyCookieConfig
	config := &GithubOauthConfig{
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		RedirectURL:  "http://localhost:9000/github/callback",
		Endpoint:     githubOAuth2.Endpoint,
	}
	//// Get the OAuth state value from the cookie
	fmt.Println("**** GetOAuthState")
	ctx, _ = GetOAuthState(cookieConfig, req, ctx)
	//oauth_state, err := GetOAuthState(cookieConfig, req)
	//if err != nil {
	//	fmt.Println("Error getting OAuth state")
	//	c.Flash.Error("Error getting OAuth state")
	//	return c.Redirect(Home.Index)
	//}
	//// Compare the state value and get the Auth token from the authorization code
	fmt.Println("**** GetToken")
	ctx = GetToken(oauth2Config, req, ctx)

	ctx = GetGithubUser(oauth2Config, req, ctx)


	//// Store the github user info in the Session cookie
	fmt.Println("**** github.UserFromContext")
	githubUser, err := github.UserFromContext(ctx)
	fmt.Println("github user", githubUser)
	fmt.Println("github user error", err)
	if err != nil {
		fmt.Println("Error getting github user from context")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return c.Redirect(Home.Index)
	}

	// TODO: Get email with GET https://api.github.com/user/emails with appropriate headers https://stackoverflow.com/questions/35373995/github-user-email-is-null-despite-useremail-scope
	emailReq, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)

	token, err := oauth2Login.TokenFromContext(ctx)
	httpClient := oauth2Config.Client(ctx, token)
	resp, err := httpClient.Do(emailReq)

	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Email response body:")
	fmt.Println(string([]byte(body)))

	type UserEmails []struct{
		Email string `json:"email"`
		Primary bool `json:"primary"`
		Verified bool `json:"verified"`
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
	fmt.Println("Primary email: ")
	fmt.Println(primaryEmail)

	// TODO: See if this user already exists in the database
	oAuthId := strconv.FormatInt(*githubUser.ID, 10)
	user, needs_creating := models.GetUserByOAuth(oAuthId, "github")
	if needs_creating {
		// Save the email address and provider in the DB
		models.SaveUser(primaryEmail, oAuthId, "github")
		user, _ = models.GetUserByOAuth(oAuthId, "github")
		fmt.Println("Saving user to DB: " + user.Email)
	} else {
		fmt.Println("Found user in DB: " + user.Email)
	}

	// Set the user ID on the user session
	c.Session["user"] = user.Id
	return c.Redirect(Home.List)
}

// Mimics gologin's StateHandler
func GetOAuthState(config gologin.CookieConfig, req *revel.Request, ctx context.Context) (context.Context, string) {
	cookie, err := req.Cookie(config.Name)
	if err == nil {
		//return cookie.GetValue(), nil
		//// add the cookie state to the ctx
		//fmt.Println("Cookie state value:", cookie.GetValue())
		//fmt.Println("state value before setting context:", ctx.Value(1))
		//fmt.Println("state value temp setting context:", oauth2Login.WithState(ctx, cookie.GetValue()).Value(1))
		ctx = oauth2Login.WithState(ctx, cookie.GetValue())
		//fmt.Println("state value after setting context:", ctx2.Value(1))
		//fmt.Println("state value after setting context:", (*ctx).val)
	} else {
		//return "", err
	}
	return ctx, cookie.GetValue()
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
		fmt.Println("authCode", authCode)// revertme
		fmt.Println("state", state)// revertme
		return "", "", errors.New("oauth2: Request missing code or state")
	}
	return authCode, state, nil
}

// Mimics gologin's oauth2.CallbackHandler
func GetToken(config *oauth2.Config, req *revel.Request, ctx context.Context) context.Context {
	authCode, state, err := parseCallback(req)
	if err != nil {
		fmt.Println("Error getting authCode from callback request:", err)
		//ctx = gologin.WithError(ctx, err)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}
	ownerState, err := oauth2Login.StateFromContext(ctx)
	if err != nil {
		fmt.Println("Error getting ownerState from context:", err)
		//ctx = gologin.WithError(ctx, err)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}
	if state != ownerState || state == "" {
		fmt.Println("oauth2 state parameter:", state) // revertme
		fmt.Println("ownerState:", ownerState) // revertme
		fmt.Println("State and ownerstate don't match:", state, ownerState)
		//ctx = gologin.WithError(ctx, ErrInvalidState)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}
	// use the authorization code to get a Token
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		fmt.Println("Error getting token from auth code:", err)
		//ctx = gologin.WithError(ctx, err)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}
	//return token
	ctx = oauth2Login.WithToken(ctx, token)
	return ctx
	//success.ServeHTTP(w, req.WithContext(ctx))
}

func GetGithubUser(config *oauth2.Config, req *revel.Request, ctx context.Context) context.Context {
	token, err := oauth2Login.TokenFromContext(ctx)
	if err != nil {
		fmt.Println("Error getting token from context:", err)
		//ctx = gologin.WithError(ctx, err)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}

	httpClient := config.Client(ctx, token)

	githubClient := go_github.NewClient(httpClient)

	user, resp, err := githubClient.Users.Get(ctx, "")
	err = validateResponse(user, resp, err)
	if err != nil {
		fmt.Println("Error validating github response:", err)
		//ctx = gologin.WithError(ctx, err)
		//failure.ServeHTTP(w, req.WithContext(ctx))
		return ctx
	}
	ctx = github.WithUser(ctx, user)
	return ctx
}

// validateResponse returns an error if the given Github user, raw
// http.Response, or error are unexpected. Returns nil if they are valid.
func validateResponse(user *go_github.User, resp *go_github.Response, err error) error {
	if err != nil || resp.StatusCode != http.StatusOK {
		return github.ErrUnableToGetGithubUser
	}
	if user == nil || user.ID == nil {
		return github.ErrUnableToGetGithubUser
	}
	return nil
}
