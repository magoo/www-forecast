package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/google/uuid"
  "github.com/aws/aws-sdk-go/aws"

)

type Scenario struct {
  Sid           string        `dynamodbav:"sid"`
  Title         string        `dynamodbav:"title"`
  Owner         string        `dynamodbav:"ownerid"`
  Hd            string        `dynamodbav:"hd"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"Options"`
  Results       []int         `dynamodbav:"results"`
  ResultIndex   int           `dynamodbav:"resultindex"`
  BrierScore    float64       `dynamodbav:"brierscore"`
  Concluded     bool          `dynamodbav:"concluded"`
  ConcludedTime string        `dynamodbav:"concludetime"`
}

func CreateScenario (title string, description string, options []string, hd string, owner string) (sid string){

  		fuuid := uuid.New()
  		item := Scenario{
  				Sid: fuuid.String(),
          Owner: owner,
          Hd: hd,
  		    Title: title,
  		    Description: description,
          Options: options,
  		}

  		PutItem(item, "scenarios-tf")

  		fmt.Println("Successfully added.")

      return fuuid.String()

}

func UpdateScenario (sid string, title string, description string, options []string, user string) {

      //Start with the key for the table
      key := map[string]*dynamodb.AttributeValue {
        "sid": {
          S: aws.String(sid),
        },
      }

      //Changing this into a list of attributes.
      o, _ := dynamodbattribute.MarshalList(options)

      expressionattrvalues:= map[string]*dynamodb.AttributeValue {
        ":t": {
          S: aws.String(title),
        },
        ":d": {
          S: aws.String(description),
        },
        ":o": {
          L: o,
        },
        ":user": {
          S: aws.String(user),
        },
      }

      updateexpression := "SET title = :t, description = :d, Options = :o"
      conditionexpression := "ownerid = :user"


  		UpdateItem(key, updateexpression, expressionattrvalues, "scenarios-tf", conditionexpression)

  		fmt.Println("Updated scenario.")
}

func ViewScenario (sid string) (s Scenario) {

  // I'll need to change this to make "secret link" work.
  result := GetPrimaryItem(sid, "sid", "scenarios-tf")

  s = Scenario{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if s.Sid == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return s

}

func ListScenarios(user string) (s []Scenario) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "scenarios-tf")

  s = []Scenario{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}

func DeleteScenario(sid string, owner string) {

  DeletePrimaryItem(sid, "sid", "scenarios-tf", "ownerid", owner)

  fmt.Println("Deleted scenario.", sid)

  DeleteScenarioForecasts(sid)

}
