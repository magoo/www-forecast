package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/magoo/www-forecast/app/controllers"
	"github.com/magoo/www-forecast/app/models"
	"golang.org/x/oauth2"

	csrf "github.com/magoo/revel-csrf"
	"github.com/revel/revel"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"

	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/sessions"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

const (
	sessionName    = "r10n-github-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "githubID"
)

// GithubOauthConfig configures the main ServeMux.
type GithubOauthConfig struct {
	GithubClientID     string
	GithubClientSecret string
}

func installHandlers() {
	revel.AddInitEventHandler(func(event revel.Event, value interface{}) (response revel.EventResponse) {
		if event == revel.ENGINE_STARTED {
			var (
				serveMux     = http.NewServeMux()
				revelHandler = revel.CurrentEngine.(*revel.GoHttpServer).Server.Handler
			)

			// (from docs) state param cookies require HTTPS by default; disable for localhost development
			var cookieConfig gologin.CookieConfig
			if revel.Config.BoolDefault("mode.dev", false) {
				cookieConfig = gologin.DebugOnlyCookieConfig
			} else {
				cookieConfig = gologin.DefaultCookieConfig
			}

			GithubOauth2Config := &oauth2.Config{
				ClientID:     controllers.GithubConfig.GithubClientID,
				ClientSecret: controllers.GithubConfig.GithubClientSecret,
				RedirectURL:  revel.Config.StringDefault("e6eDomain", "https://localhost:9000") + "/github/callback",
				Endpoint:     controllers.GithubEndpoint,
			}
			TwitterConfig := &oauth1.Config{
				ConsumerKey:    os.Getenv("TWITTER_CLIENT_ID"),
				ConsumerSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
				CallbackURL:    revel.Config.StringDefault("e6eDomain", "https://localhost:9000") + "/twitter/callback",
				Endpoint:       twitterOAuth1.AuthorizeEndpoint,
			}

			serveMux.Handle("/github/login", github.StateHandler(cookieConfig, github.LoginHandler(GithubOauth2Config, LoggingErrorHandler)))
			serveMux.Handle("/twitter/login", twitter.LoginHandler(TwitterConfig, LoggingErrorHandler))

			serveMux.Handle("/", revelHandler) // Should this be "*" or something?
			revel.CurrentEngine.(*revel.GoHttpServer).Server.Handler = serveMux
		}
		return 0
	})
}

func errorLogger(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := gologin.ErrorFromContext(ctx)
	if err != nil {
		fmt.Println(err)
		fmt.Println(req)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

var LoggingErrorHandler = http.HandlerFunc(errorLogger)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// issueSession issues a cookie session after successful Github login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		fmt.Println("context:", ctx)
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID

		// TODO: Figure out how to use Revel Controller session variables for c.Session["user"] = user.Id

		fmt.Println("github user:", *githubUser)

		//email := ""
		//if githubUser.Email != nil {
		//	email = *githubUser.Email
		//}

		models.CreateUser("test@example.com", strconv.FormatInt(*githubUser.ID, 10), "github")

		session.Save(w)
		http.Redirect(w, req, "/list", http.StatusFound) // TODO: Make this URL something that we want
	}
	return http.HandlerFunc(fn)
}

func init() {
	revel.OnAppStart(installHandlers)
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		//added for CSRF protection
		csrf.CSRFFilter,        // CSRF prevention.
		revel.ValidationFilter, // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,       // Resolve the requested language
		HeaderFilter,           // Add some security based headers
		AuthFilter,
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.ActionInvoker,     // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)

	revel.TemplateFuncs["animalFaceEmoji"] = func(index int) string {
		var animalFaceEmojis = [...]rune{
			'🐵', '🐶', '🐺', '🦊', '🐱',
			'🦁', '🐯', '🐴', '🦄', '🐮',
			'🐷', '🐗', '🐭', '🐹', '🐰',
			'🐻', '🐨', '🐼', '🐔', '🐤',
			'🐦', '🐧', '🐸', '🐲', '🐙',
		}

		// mod to allow wrap around if the index is out of bounds
		return string(animalFaceEmojis[index%25])
	}

	revel.TemplateFuncs["brierRound"] = func(bs float64) float64 {
		return models.RoundPlus(bs, 3)
	}

	revel.TemplateFuncs["brierColor"] = func(bs float64) string {

		switch {
		case (bs < .10):
			return "#25A400"
		case bs > .10 && bs < .15:
			return "#25A400"
		case bs > .15 && bs < .20:
			return "#25A400"
		case bs > .20 && bs < .25:
			return "#D2FDC5"
		case bs > .25 && bs < .30:
			return "#FDCFC5"
		case bs > .30 && bs < .45:
			return "#FEBBAD"
		case bs > .45 && bs < .55:
			return "#FE9B86"
		case bs > .65 && bs < .75:
			return "#FF7355"
		case bs > .75 && bs < .85:
			return "#FF5733"
		case bs > .85 && bs < 1:
			return "#A8432E"
		default:
			return "black"
		}

	}

	revel.TemplateFuncs["indexToCharacter"] = func(index int) string {
		character := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		// mod to allow wrap around if the index is out of bounds
		return string(character[index%26])
	}

	if cfgPaths, ok := os.LookupEnv("E6E_CONFIG_PATHS"); ok {
		revel.ConfPaths = strings.Split(cfgPaths, ",")
	}
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	if revel.RunMode == "prod" {
		c.Response.Out.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

var AuthFilter = func(c *revel.Controller, fc []revel.Filter) {

	gc := revel.Config.StringDefault("e6e.google_client", "empty")
	c.ViewArgs["gc"] = gc

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
