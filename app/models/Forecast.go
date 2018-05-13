package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/docker/docker/pkg/namesgenerator"
  //"os"
  "time"

)

type Forecast struct {
  Hd            string        `dynamodbav:"hd"`
  Sid           string        `dynamodbav:"sid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"ownerid"`
  Forecasts     []int         `dynamodbav:"forecasts"`
  UserAlias     string        `dynamodbav:"useralias"`
}

func CreateForecast (u string, f []int, sid string, hd string) {
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

      t := time.Now()

  		item := Forecast{
          Hd: hd,
          Sid: sid,
          Date: t.String(),
  		    User: u,
          UserAlias: namesgenerator.GetRandomName(0),
  		    Forecasts: f,
  		}

  		PutItem(item, "forecasts-tf")

}

func ViewScenarioResults (sid string) (c []Forecast) {
  //Need to do a HD check here to prevent IDOR.

    result := GetPrimaryIndexItem(sid, "sid", "sid-index", "forecasts-tf")

    c = []Forecast{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &c)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return c
}

func ViewUserScenarioResults (uid string, sid string) (userForecast Forecast) {
  results := ViewScenarioResults(sid)
  for _, result := range results {
    if result.User == uid {
      userForecast = result
    }
  }

  return userForecast
}

func DeleteScenarioForecasts(sid string) {

    fs := ViewScenarioResults(sid)


    for _, v  := range fs {
      fmt.Println("Deleting: ", v.Sid, v.User)
      DeleteCompositeIndexItem(v.Sid, v.User, "sid", "user", "forecasts-tf")
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
