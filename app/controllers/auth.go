package controllers

import (
	"github.com/revel/revel"
	"google.golang.org/api/oauth2/v2"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"

)

type Oauth2 struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Exp           string `json:"exp"`
	Iss           string `json:"iss"`
	Jti           string `json:"jti"`
	Iat           string `json:"iat"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
}

type Auth struct {
	*revel.Controller
}

var httpClient = &http.Client{}

func (c Auth) Index() revel.Result {

	c.Validation.Required(c.Session["user"])
	c.Validation.Required(c.Session["hd"])

	if c.Validation.HasErrors() {
		//fmt.Println("We had a validation error")
		return c.Render(Auth.Index)
	}

	return c.Redirect(List.Index)
}

func (c Auth) GoogleToken(idtoken string) revel.Result {
	//fmt.Println("Received: ", idtoken)

	ti, hd, err := verifyIdToken(idtoken)

	//fmt.Println("verified email: ", ti.Email)
	//fmt.Println("verified domain: ", hd)
	//fmt.Printf("%+v\n", ti)
	if (err != nil){
		panic("Can't verify ID token from Google: " + err.Error())
	}

	c.Session["user"] = ti.Email

	// Set the hosted domain. This is our core privacy barrier.
	// If not a gsuite customer, it's public.
	// If gsuite customer, it's for the "hosted domain" only.
	// Cheap and simple privacy for the time being.
	// We change empty strings to "public" because dynamo doesn't like empty
	// strings downstream.
	if (hd == "") {
		c.Session["hd"] = "public"
	} else {
		c.Session["hd"] = hd
	}

	return c.Render()
}

func verifyIdToken(idToken string) (*oauth2.Tokeninfo, string, error) {
    oauth2Service, err := oauth2.New(httpClient)
    tokenInfoCall := oauth2Service.Tokeninfo()
    tokenInfoCall.IdToken(idToken)
    tokenInfo, err := tokenInfoCall.Do()
    if err != nil {
        panic("Can't verify hd claim from Google: " + err.Error())
    }

		hd := getHd(idToken)
    return tokenInfo, hd, nil
}

//This function exists because of a shortcoming in the Google oauth2 api. `hd` is not returned in verifyIdToken(). Causes a second request to google on login.
func getHd(idToken string)(hd string){
	client := &http.Client{}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = true
	client.Transport = transport
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken, nil)
  if err != nil {
			 fmt.Println(err)
  }
	res, err:= client.Do(req)
	defer res.Body.Close()

	fmt.Println(res)

	if err != nil {
		fmt.Println("Request failed.")
    panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Reading body failed.")
    panic(err.Error())
	}

	o := Oauth2{}

	err = json.Unmarshal(body, &o)

	if (err != nil) {
		fmt.Println("Unmarshal failed.")
		fmt.Println("error:", err)
	}
	return o.Hd

}
