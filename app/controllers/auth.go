package controllers

import (
	"github.com/revel/revel"
	"google.golang.org/api/oauth2/v2"
	"net/http"

)

type Auth struct {
	*revel.Controller
}

var httpClient = &http.Client{}

func (c Auth) Index() revel.Result {
	return c.Render()

}

func (c Auth) GoogleToken(idtoken string) revel.Result {
	fmt.Println("Received: ", idtoken)

	ti, err := verifyIDToken(idtoken)


	return c.Render()
}

func verifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
    oauth2Service, err := oauth2.New(httpClient)
    tokenInfoCall := oauth2Service.Tokeninfo()
    tokenInfoCall.IdToken(idToken)
    tokenInfo, err := tokenInfoCall.Do()
    if err != nil {
        return nil, err
    }
    return tokenInfo, nil
}
