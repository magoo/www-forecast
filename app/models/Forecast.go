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
  Hd            string        `dynamodbav:"hd"`
  Sid           string        `dynamodbav:"sid"`
  Date          string        `dynamodbav:"date"`
  User          string        `dynamodbav:"user"`
  Forecasts     []int         `dynamodbav:"forecasts"`
}

func CreateForecast (u string, f []int, sid string, hd string) (fid string){
      //Must do a permission check in the future to prevent crossover forecasts. Tock day.
      //Must do a check to make sure the array of values is equal to the array of options in the sid.

  		fuuid := uuid.New()
      t := time.Now()
      fmt.Println(sid)
  		item := Forecast{
  				Fid: fuuid.String(),
          Hd: hd,
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

func ViewScenarioResults (sid string, hd string) (c []Forecast) {
  //Need to do a HD check here to prevent IDOR.

    input := &dynamodb.QueryInput{
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":v1": {
                S: aws.String(sid),
            },
            ":v2": {
                S: aws.String(hd),
            },
        },
        KeyConditionExpression: aws.String("sid = :v1 AND hd = :v2"),
        IndexName:              aws.String("sid-hd-index"),
        TableName:              aws.String("forecasts"),
    }

    result, err := Svc.Query(input)
    if err != nil {
            fmt.Println("Error viewing scenario results: " , err.Error())
    }

    c = []Forecast{}

    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &c)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return c
}

func DeleteScenarioForecasts(sid string, hd string) {

    fs := ViewScenarioResults(sid, hd)


    for _, v  := range fs {
      fmt.Println("Deleting: ", v.Fid)
      deleteRequest := &dynamodb.DeleteItemInput{
          Key: map[string]*dynamodb.AttributeValue{
              "fid": {
                  S: aws.String(v.Fid),
              },
              "sid": {
                  S: aws.String(v.Sid),
              },
          },
          TableName: aws.String("forecasts"),
      }

      _, err := Svc.DeleteItem(deleteRequest)



    if err != nil {
              fmt.Println("Got error calling DeleteItem")
              fmt.Println(err.Error())
              return
          }
    }


    fmt.Println("Deleted forecasts associated with scenario.")

}
