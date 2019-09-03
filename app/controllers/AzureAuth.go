package controllers

import (
	"context"
	"crypto/rand"
	"errors"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/gologin"
	oauth2Login "github.com/dghubble/gologin/oauth2"

	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"

	"os"

	"golang.org/x/oauth2"
)

type AzureAuth struct {
	*revel.Controller
}

type AzureUser struct {
	Id                string `json:"id"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

type AzureOauthConfig struct {
	AzureClientID     string
	AzureObjectID     string
	AzureClientSecret string
}

var (
	azureEndpoint = oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}
	config = &AzureOauthConfig{
		AzureClientID:     os.Getenv("AZURE_CLIENT_ID"),
		AzureObjectID:     os.Getenv("AZURE_OBJECT_ID"),
		AzureClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
	}
)

var (
	ErrUnableToGetAzureUser = errors.New("azure: unable to get Azure User")
)

// Returns a base64 encoded random 32 byte string.
func randomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (c AzureAuth) Login() revel.Result {
	var cookieConfig gologin.CookieConfig
	if revel.Config.BoolDefault("mode.dev", false) {
		cookieConfig = gologin.DebugOnlyCookieConfig
	} else {
		cookieConfig = gologin.DefaultCookieConfig
	}

	ctx := c.Request.Context()
	cookie, err := c.Request.Cookie(cookieConfig.Name)
	state := ""
	if err == nil {
		// add the cookie state to the ctx
		state = cookie.GetValue()
		ctx = oauth2Login.WithState(ctx, state)
	} else {
		// add Cookie with a random state
		state = randomState()
		c.SetCookie(&http.Cookie{Name: cookieConfig.Name, Value: state})
		ctx = oauth2Login.WithState(ctx, state)
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.AzureClientID,
		ClientSecret: config.AzureClientSecret,
		RedirectURL:  revel.Config.StringDefault("e6eDomain", "https://localhost:9000") + "/azure/callback",
		Endpoint:     azureEndpoint,
		Scopes: []string{
			"https://graph.microsoft.com/User.Read",
		},
	}
	url := oauth2Config.AuthCodeURL(state)
	c.Response.ContentType = "application/x-www-form-urlencoded"
	return c.Redirect(url)

}

func GetAzureToken(config *oauth2.Config, req *revel.Request, ctx context.Context, ownerState string) *oauth2.Token {
	authCode, state, err := parseAzureCallback(req)
	if err != nil {
		fmt.Println("Error getting authCode from callback request:", err)
		return nil
	}
	if state != ownerState || state == "" {
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

func parseAzureCallback(req *revel.Request) (authCode, state string, err error) {
	err = req.ParseForm()
	if err != nil {
		return "", "", err
	}
	authCode = req.Form.Get("code")
	state = req.Form.Get("state")
	if authCode == "" || state == "" {
		return "", "", errors.New("oauth2: Request missing code or state")
	}
	return authCode, state, nil
}

func GetAzureOAuthState(config gologin.CookieConfig, req *revel.Request, ctx context.Context) string {
	cookie, err := req.Cookie(config.Name)
	if err == nil {
		return cookie.GetValue()
	} else {
		return ""
	}
}

func (c AzureAuth) Callback() revel.Result {
	ctx := c.Request.Context()

	req := c.Request
	// TODO: Include MIT license with this code taken from the library
	// TODO: Dry out with init.go
	cookieConfig := gologin.DebugOnlyCookieConfig
	//// Get the OAuth state value from the cookie
	ownerState := GetAzureOAuthState(cookieConfig, req, ctx)

	oauth2Config := &oauth2.Config{
		ClientID:     config.AzureClientID,
		ClientSecret: config.AzureClientSecret,
		RedirectURL:  revel.Config.StringDefault("e6eDomain", "https://localhost:9000") + "/azure/callback",
		Endpoint:     azureEndpoint,
		Scopes: []string{
			"https://graph.microsoft.com/User.Read",
		},
	}

	//// Compare the state value and get the Auth token from the authorization code
	token := GetAzureToken(oauth2Config, req, ctx, ownerState)
	fmt.Println("token:", token)

	azureUser, err := GetAzureUser(oauth2Config, ctx, token)
	if err != nil {
		return c.RenderError(err)
	}

	//// See if this user already exists in the database
	oAuthId := azureUser.Id
	user, needs_creating := models.GetUserByOAuth(oAuthId, "azure")

	// Azure users don't always have emails, but their principal names will exist and be emails in at least some of these cases
	email := ""
	if azureUser.Mail != "" {
		email = azureUser.Mail
	} else if azureUser.UserPrincipalName != "" {
		email = azureUser.UserPrincipalName
	}
	if needs_creating {
		// Save the email address and provider in the DB
		models.CreateUser(email, oAuthId, "azure")
		user, _ = models.GetUserByOAuth(oAuthId, "azure")
		fmt.Println("Saving user to DB: " + email)
	} else {
		fmt.Println("Found user in DB: " + email)
	}

	// Set the user ID on the user session
	c.Session["user"] = user.OauthProvider + ":" + user.OauthID
	return c.Redirect(Home.List)
}

func GetAzureUser(config *oauth2.Config, ctx context.Context, token *oauth2.Token) (*AzureUser, error) {
	httpClient := config.Client(ctx, token)

	// Microsoft graph endpoint that gives information about the user we're authenticating with
	userReq, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)

	resp, err := httpClient.Do(userReq)
	if err != nil {
		return nil, err
	}
	user := new(AzureUser)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal([]byte(body), &user)

	return user, nil
}
