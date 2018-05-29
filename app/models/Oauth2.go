package models

import (
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

var httpClient = &http.Client{}

func VerifyIdToken(idToken string) (*oauth2.Tokeninfo, string, error) {
    oauth2Service, err := oauth2.New(httpClient)

		if err != nil {
			fmt.Println("Cannot create HTTP client: " + err.Error())
		}

    tokenInfoCall := oauth2Service.Tokeninfo()
    tokenInfoCall.IdToken(idToken)
    tokenInfo, err := tokenInfoCall.Do()

		if err != nil {
			fmt.Println("Cannot get OAuth data: " + err.Error())
		}

		hd, err := getHd(idToken)

    return tokenInfo, hd, err
}

//This function exists because of a shortcoming in the Google oauth2 api. `hd` is not returned in verifyIdToken(). Causes a second request to google on login.
func getHd(idToken string) (hd string, err error) {
	client := &http.Client{}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = true
	client.Transport = transport
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken, nil)
  if err != nil {
			 fmt.Println("HTTP error crafting Google API Token Request" + err.Error())
  }
	res, err:= client.Do(req)
	defer res.Body.Close()

	if err != nil {
			 fmt.Println("HTTP error getting Google API ID Token" + err.Error())
	}

	//fmt.Println(res)

	body, err := ioutil.ReadAll(res.Body)

	o := Oauth2{}

	err = json.Unmarshal(body, &o)
	if err != nil {
			 fmt.Println("Error with JSON unmarshal" + err.Error())
  }
	return o.Hd, err

}
