package models

import (
  //"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/docker/docker/pkg/namesgenerator"
  //"os"
  "time"
  "fmt"

)


type Sort struct {
  Answer
  Options       []int      `dynamodbav:"options"`
}

func CreateSort (u string, options []int, rid string, hd string) {
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

      t := time.Now()

  		item := Sort{
        Answer: Answer{
            Hd: hd,
            Id: rid,
            Date: t.String(),
            OwnerID: u,
            UserAlias: namesgenerator.GetRandomName(0),
          },
          Options: options,
  		}

  		err := PutItem(item, answerTable)

      if err != nil {
        fmt.Println("Error writing to db.")
      } else {
        fmt.Println("Successfully added.")
      }

}
