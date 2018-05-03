package models

import (
  //"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/docker/docker/pkg/namesgenerator"
  //"os"
  "time"

)

type Sort struct {
  Hd            string        `dynamodbav:"hd"`
  Rid           string        `dynamodbav:"rid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"ownerid"`
  Options       []int      `dynamodbav:"options"`
  UserAlias     string        `dynamodbav:"useralias"`
}

func CreateSort (u string, options []int, rid string, hd string) {
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

      t := time.Now()

  		item := Sort{
          Hd: hd,
          Rid: rid,
          Date: t.String(),
          Options: options,
  		    User: u,
          UserAlias: namesgenerator.GetRandomName(0),
  		}

  		PutItem(item, "sorts-tf")

}
