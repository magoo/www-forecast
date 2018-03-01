package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/google/uuid"
)

type Scenario struct {
  Sid           string        `dynamodbav:"sid"`
  Title         string        `dynamodbav:"title"`
  Hd            string        `dynamodbav:"hd"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"Options"`
  Results       []int         `dynamodbav:"results"`
  ResultIndex   int           `dynamodbav:"resultindex"`
  BrierScore    float64       `dynamodbav:"brierscore"`
  Concluded     bool          `dynamodbav:"concluded"`
  ConcludedTime  string        `dynamodbav:"concludetime"`
}

func CreateScenario (title string, description string, options []string, hd string) (sid string){

  		fuuid := uuid.New()
  		item := Scenario{
  				Sid: fuuid.String(),
          Hd: hd,
  		    Title: title,
  		    Description: description,
          Options: options,
  		}

  		PutItem(item, "scenarios")

  		fmt.Println("Successfully added.")

      return fuuid.String()

}

func ViewScenario (sid string, hd string) (s Scenario) {

  result := GetCompositeKeyItem(sid, hd, "sid", "hd", "scenarios")

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

func ListScenarios(hd string) (s []Scenario) {
  //This must respects "hd" privacy. Only return results from the "Hosted Domain" in Google.

  result := GetPrimaryIndexItem(hd, "hd", "hd-index", "scenarios")

  s = []Scenario{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}

func DeleteScenario(sid string, hd string) {

  DeleteCompositeIndexItem(sid, hd, "sid", "hd", "scenarios")

  fmt.Println("Deleted scenario.")

  DeleteScenarioForecasts(sid, hd)

}
