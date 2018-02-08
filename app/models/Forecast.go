package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  "os"
  "github.com/google/uuid"
  "time"

)

type Forecast struct {
  Fid           string        `dynamodbav:"fid"`
  Sid           string        `dynamodbav:"sid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"user"`
  Forecasts     []int         `dynamodbav:"forecasts"`
}

func CreateForecast (u string, f []int, sid string) (fid string){
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

  		fuuid := uuid.New()
      t := time.Now()
      fmt.Println(sid)
  		item := Forecast{
  				Fid: fuuid.String(),
          Sid: sid,
          Date: t.String(),
  		    User: u,
  		    Forecasts: f,
  		}

  		av, err := dynamodbattribute.MarshalMap(item)

  		if err != nil {
  			fmt.Println("Got error calling MarshalMap:")
  			fmt.Println(err.Error())
  			os.Exit(1)
  		}

  		input := &dynamodb.PutItemInput{
  	    Item: av,
  	    TableName: aws.String("forecasts"),
  		}

  		_, err = Svc.PutItem(input)

  		if err != nil {
  	    fmt.Println("Got error calling PutItem:")
  	    fmt.Println(err.Error())
  	    os.Exit(1)
  		}

  		fmt.Println("Successfully added.")

      //Return the cast id
      return fuuid.String()

}

func ViewScenarioResults (sid string) (c []Forecast) {
  //Need to do a HD check here to prevent IDOR.

    input := &dynamodb.QueryInput{
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":v1": {
                S: aws.String(sid),
            },
        },
        KeyConditionExpression: aws.String("sid = :v1"),
        IndexName:              aws.String("sid-index"),
        TableName:              aws.String("forecasts"),
    }

    result, err := Svc.Query(input)
    if err != nil {
            fmt.Println(err.Error())
    }

    c = []Forecast{}

    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &c)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return c
}
