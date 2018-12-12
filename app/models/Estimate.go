package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  //"os"
  "github.com/google/uuid"
  "time"
  "strconv"
  "errors"

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
                  Date: t.Format("2006-01-02"),
                  OwnerID: owner,
                  Hd: hd,
          		    Title: title,
          		    Description: description,
                  Records: []string{t.Format("2006-01-02") + ": Created.", },
                  URL: "estimate/" + euuid.String(),
                  Type: "Estimate",
                },
          Unit: unit,
  		}

  		err := PutItem(item, questionTable)

      if err != nil {
        fmt.Println("Error writing to db.")
      } else {
        fmt.Println("Successfully added.")
      }

      return euuid.String()
}

func (e Estimate) GetURL() (url string) {

  return "/view/estimate/" + e.Id

}

func (e Estimate) AddRecord(user string) (err error){

  er := ViewEstimateResults(e.Question.Id)

  if len(er) < 1 {
    return errors.New("No results.")
  }

  emin, emax := GetAverageRange(er)

  t := time.Now()

  record := t.Format("2006-01-02") + ": Results recorded. Min: " + strconv.FormatFloat(emin, 'f', -1, 64) + " Max: " + strconv.FormatFloat(emax, 'f', -1, 64)

  //Primary key for update query
  key := map[string]*dynamodb.AttributeValue {
    "id": {
      S: aws.String(e.Question.Id),
    },
  }

  item := map[string]*dynamodb.AttributeValue {
    ":r": {
        SS: []*string{
          aws.String(record),
          },
        },
    ":user": {
      S: aws.String(user),
    },
  }

  //av, err := dynamodbattribute.MarshalMap(item)

  UpdateItem(key, "ADD records :r", item, questionTable, "ownerid = :user")

  return err

}

func GetAverageRange(er []Range) (avgmin float64, avgmax float64){

	size := len(er)
	var sum float64 = 0

	for _, v := range er {
		sum += v.Minimum
	}

	avgmin = sum / float64(size)

	//reset
	sum = 0

	for _, v := range er {
		sum += v.Maximum
	}

	avgmax = sum / float64(size)

	return
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

  UpdateItem(key, updateexpression, expressionattrvalues, questionTable, conditionexpression)

}

func GetEstimate (eid string) (e Estimate) {

  result := GetPrimaryItem(eid, "id", questionTable)

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

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", questionTable)

  e = []Estimate{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &e)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return e

}

func DeleteEstimate(eid string, owner string) {

  DeletePrimaryItem(eid, "id", questionTable, "ownerid", owner)

  fmt.Println("Deleted estimate.", eid)

  DeleteEstimateRanges(eid)

}

func DeleteEstimateRanges(eid string) {

    er := ViewEstimateResults(eid)


    for _, v  := range er {
      fmt.Println("Deleting: ", v.Answer.Id, v.Answer.OwnerID)
      DeleteCompositeIndexItem(v.Answer.Id, v.Answer.OwnerID, "id", "ownerid", "answerTable")
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
