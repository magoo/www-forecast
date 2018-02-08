package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  "os"
  "github.com/google/uuid"
)

type Scenario struct {
  Sid           string        `dynamodbav:"sid"`
  Title         string        `dynamodbav:"title"`
  Hd            string        `dynamodbav:"hd"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"Options"`
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

  		av, err := dynamodbattribute.MarshalMap(item)

  		if err != nil {
  			fmt.Println("Got error calling MarshalMap:")
  			fmt.Println(err.Error())
  			os.Exit(1)
  		}

  		input := &dynamodb.PutItemInput{
  	    Item: av,
  	    TableName: aws.String(dbname),
  		}

  		_, err = Svc.PutItem(input)

  		if err != nil {
  	    fmt.Println("Got error calling PutItem:")
  	    fmt.Println(err.Error())
  	    os.Exit(1)
  		}

  		fmt.Println("Successfully added.")

      return fuuid.String()

}

func ViewScenario (sid string, hd string) (s Scenario) {
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        "sid": {
            S: aws.String(sid),
        },
        "hd": {
            S: aws.String(hd),
        },
    },
    TableName: aws.String(dbname),
  }

  result, err := Svc.GetItem(input)
  if err != nil {
          fmt.Println(err.Error())
  }

  s = Scenario{}

  err = dynamodbattribute.UnmarshalMap(result.Item, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if s.Sid == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return s

}


func ListScenarios (hd string) (s []Scenario) {
  //This must respects "hd" privacy. Only return results from the "Hosted Domain" in Google.

  input := &dynamodb.QueryInput{
      ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
          ":v1": {
              S: aws.String(hd),
          },
      },
      KeyConditionExpression: aws.String("hd = :v1"),
      IndexName:              aws.String("hd-index"),
      TableName:              aws.String(dbname),
  }

  result, err := Svc.Query(input)
  if err != nil {
          fmt.Println(err.Error())
  }

  s = []Scenario{}

  err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}
