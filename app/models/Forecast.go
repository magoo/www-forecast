package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  "os"
  "github.com/google/uuid"
)

type Forecast struct {
  Fid           string        `dynamodbav:"Fid"`
  Title         string        `dynamodbav:"title"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"Options"`
}

func CreateForecast (title string, description string, options []string) (fid string){

  		fuuid := uuid.New()

  		item := Forecast{
  				Fid: fuuid.String(),
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
  	    TableName: aws.String("testing_table"),
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

func ViewForecast (fid string) (f Forecast) {
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        "Fid": {
            S: aws.String(fid),
        },
    },
    TableName: aws.String("testing_table"),
  }

  result, err := Svc.GetItem(input)
  if err != nil {
          fmt.Println(err.Error())
  }

  f = Forecast{}

  err = dynamodbattribute.UnmarshalMap(result.Item, &f)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if f.Fid == "" {
      fmt.Println("Could not find that forecast.")
      return
  }

  return f

}

func ListForecasts () (f []Forecast) {
  //This will eventually break when scans are greater than 1mb

  input := &dynamodb.ScanInput{
    TableName:            aws.String("testing_table"),
  }

  result, err := Svc.Scan(input)
  if err != nil {
          fmt.Println(err.Error())
  }

  f = []Forecast{}

  err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &f)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return f

}
