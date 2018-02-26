package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  //"os"
  "time"

)

type Forecast struct {
  Hd            string        `dynamodbav:"hd"`
  Sid           string        `dynamodbav:"sid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"user"`
  Forecasts     []int         `dynamodbav:"forecasts"`
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
  		    Forecasts: f,
  		}

  		PutItem(item, "forecasts")

}

func ViewScenarioResults (sid string, hd string) (c []Forecast) {
  //Need to do a HD check here to prevent IDOR.

    result := GetCompositeIndexItem(sid, hd, "sid", "hd", "sid-hd-index", "forecasts")

    c = []Forecast{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &c)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return c
}

func DeleteScenarioForecasts(sid string, hd string) {

    fs := ViewScenarioResults(sid, hd)


    for _, v  := range fs {
      fmt.Println("Deleting: ", v.Sid, v.User)
      DeleteCompositeIndexItem(v.Sid, v.User, "sid", "user", "forecasts")
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
