package models

import (

  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "fmt"

)

type Question struct {
  Id            string        `dynamodbav:"id"` //Uniquely identify the question
  OwnerID       string        `dynamodbav:"ownerid"` // Owner of the question, is moderator.
  Date          string        `dynamodbav:"date"`
  Hd            string        `dynamodbav:"hd"` // Owning organization (for larger group visibility)
  Title         string        `dynamodbav:"title"`
  Description   string        `dynamodbav:"description"`
  BrierScore    float64       `dynamodbav:"brierscore"` // Any rolling Brier score we are calculating
  Concluded     bool          `dynamodbav:"concluded"` // Has this scenario shut down?
  ConcludedTime string        `dynamodbav:"concludetime"` //If so, when?
  Records       []string      `dynamodbav:"records,stringset"`    // Audit records on the scenario.
  URL           string        `dynamodbav:"url"`
  Type          string        `dynamodbav:"type"`
}

type Answer struct {
  Id            string        `dynamodbav:"id"` // The question this answers
  OwnerID       string        `dynamodbav:"ownerid"` // Owner of the question, is moderator.
  Hd            string        `dynamodbav:"hd"`
  Date          string        `dynamodbav:"date"`
  UserAlias     string        `dynamodbav:"useralias"` // The users fake name
  URL           string        `dynamodbav:"url"`
  Title         string        `dynamodbav:"title"`
  Description   string        `dynamodbav:"description"`
}

func ListQuestions(user string) (s []Question) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "questions-tf")

  s = []Question{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}


func ListAnswers(user string) (s []Question) {
  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "answers-tf")

  s = []Question{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}
