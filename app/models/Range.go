package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/docker/docker/pkg/namesgenerator"
  //"os"
  "time"

)

type Range struct {
  Hd            string        `dynamodbav:"hd"`
  Eid           string        `dynamodbav:"eid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"ownerid"`
  Minimum       float64       `dynamodbav:"minimum"`
  Maximum       float64       `dynamodbav:"maximum"`
  UserAlias     string        `dynamodbav:"useralias"`
}

func CreateRange (u string, min float64, max float64, eid string, hd string) {
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

      t := time.Now()

  		item := Range{
          Hd: hd,
          Eid: eid,
          Date: t.String(),
  		    User: u,
          UserAlias: namesgenerator.GetRandomName(0),
  		    Minimum: min,
          Maximum: max,
  		}

  		PutItem(item, "ranges-tf")

}

func ViewEstimateResults (eid string) (rs []Range) {
  //Need to do a HD check here to prevent IDOR.

    result := GetPrimaryIndexItem(eid, "eid", "eid-index", "ranges-tf")

    rs = []Range{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &rs)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return rs
}
