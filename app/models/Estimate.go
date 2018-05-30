package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  //"os"
  "github.com/google/uuid"
  "time"

)

type Estimate struct {
  Question
  AvgMinimum    float64       `dynamodbav:"minimum"`
  AvgMaximum    float64       `dynamodbav:"maximum"`
  Actual        float64       `dynamodbav:"actual"`
  Unit          string        `dynamodbav:"unit"`
}

func CreateEstimate (title string, description string, unit string, hd string, owner string) (eid string){

      t := time.Now()

  		euuid := uuid.New()
  		item := Estimate{
        Question: Question{
          				Id: euuid.String(),
                  OwnerID: owner,
                  Hd: hd,
          		    Title: title,
          		    Description: description,
                  Records: []string{t.Format("2006-01-02") + ": Created.", },
                  URL: "estimate/" + euuid.String(),
                },
          Unit: unit,
  		}

  		PutItem(item, "questions-tf")
      fmt.Println(unit)

      return euuid.String()
}

func (e Estimate) GetURL() (url string) {

  return "/view/estimate/" + e.Id

}

func UpdateEstimate (eid string, title string, description string, unit string, user string) {

  //Primary key for update query
  key := map[string]*dynamodb.AttributeValue {
    "id": {
      S: aws.String(eid),
    },
  }

  expressionattrvalues:= map[string]*dynamodb.AttributeValue {
    ":t": {
      S: aws.String(title),
    },
    ":d": {
      S: aws.String(description),
    },
    ":unit": {
      S: aws.String(unit),
    },
    ":user": {
      S: aws.String(user),
    },
  }

  updateexpression := "SET title = :t, description = :d, unitname = :unit"
  conditionexpression := "ownerid = :user"

  UpdateItem(key, updateexpression, expressionattrvalues, "questions-tf", conditionexpression)

}

func GetEstimate (eid string) (e Estimate) {

  result := GetPrimaryItem(eid, "id", "questions-tf")

  e = Estimate{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &e)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if e.Question.Id == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return e

}

func ListEstimates(user string) (e []Estimate) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "questions-tf")

  e = []Estimate{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &e)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return e

}

func DeleteEstimate(eid string, owner string) {

  DeletePrimaryItem(eid, "id", "questions-tf", "ownerid", owner)

  fmt.Println("Deleted estimate.", eid)

  DeleteEstimateRanges(eid)

}

func DeleteEstimateRanges(eid string) {

    er := ViewEstimateResults(eid)


    for _, v  := range er {
      fmt.Println("Deleting: ", v.Answer.Id, v.Answer.OwnerID)
      DeleteCompositeIndexItem(v.Answer.Id, v.Answer.OwnerID, "eid", "user", "answers-tf")
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
