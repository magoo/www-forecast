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
  Hd            string        `dynamodbav:"hd"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"Options"`
}

func CreateForecast (title string, description string, options []string, hd string) (fid string){

  		fuuid := uuid.New()
  		item := Forecast{
  				Fid: fuuid.String(),
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

func ViewForecast (fid string, hd string) (f Forecast) {
  input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        "Fid": {
            S: aws.String(fid),
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


func ListForecasts (hd string) (f []Forecast) {
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

  f = []Forecast{}

  err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &f)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return f

}
