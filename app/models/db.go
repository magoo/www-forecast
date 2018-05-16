package models

import (
		"github.com/aws/aws-sdk-go/aws"
	  "github.com/aws/aws-sdk-go/aws/session"
	  "github.com/aws/aws-sdk-go/service/dynamodb"
		"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
		"fmt"
		"os"
)

var sess, _ = session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	},
)

var dbname = "scenarios-tf"

var Svc = dynamodb.New(sess)

func PutItem(item interface{}, table string) (err error) {

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		fmt.Println("Got error calling MarshalMap: ")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(table),
	}

	_, err = Svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem: ")
		fmt.Println(err.Error())

	}
	return
}

func UpdateItem(key map[string]*dynamodb.AttributeValue, updateexpression string, expressionattrvalues map[string]*dynamodb.AttributeValue, table string, conditionexpression string ) (err error) {

	input := &dynamodb.UpdateItemInput{
    ExpressionAttributeValues: expressionattrvalues,
    Key: key,
    TableName:        aws.String(table),
    UpdateExpression: aws.String(updateexpression),
		ConditionExpression: aws.String(conditionexpression),
	}

	if err != nil {
		fmt.Println("Got error calling MarshalMap: ")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = Svc.UpdateItem(input)

	if err != nil {
		fmt.Println("Got error calling UpdateItem: ")
		fmt.Println(err.Error())

	}
	return
}

func GetPrimaryIndexItem(primaryValue string, primary string, index string, table string) (result *dynamodb.QueryOutput){
	input := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":v1": {
							S: aws.String(primaryValue),
					},
			},
			KeyConditionExpression: aws.String(primary + " = :v1"),
			IndexName:              aws.String(index),
			TableName:              aws.String(table),
	}

	result, err := Svc.Query(input)
	if err != nil {
					fmt.Println(err.Error())
	}

	return

}

func GetPrimaryItem(primaryValue string, primary string, table string) (result *dynamodb.GetItemOutput){
	input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        primary: {
            S: aws.String(primaryValue),
        },
    },
    TableName: aws.String(table),
  }

  result, err := Svc.GetItem(input)
  if err != nil {
          fmt.Println(err.Error())
  }

	return
}

func GetCompositeKeyItem(primaryValue string, sortValue string, primary string, sort string, table string) (result *dynamodb.GetItemOutput) {

	input := &dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
        primary: {
            S: aws.String(primaryValue),
        },
        sort: {
            S: aws.String(sortValue),
        },
    },
    TableName: aws.String(table),
  }

  result, err := Svc.GetItem(input)
  if err != nil {
          fmt.Println(err.Error())
  }

	return

}

func GetCompositeIndexItem(primaryValue string, sortValue string, primary string, sort string, index string, table string) (result *dynamodb.QueryOutput) {

	input := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":v1": {
							S: aws.String(primaryValue),
					},
					":v2": {
							S: aws.String(sortValue),
					},
			},
			KeyConditionExpression: aws.String(primary + " = :v1 AND " + sort + " = :v2"),
			IndexName:              aws.String(index),
			TableName:              aws.String(table),
	}

	result, err := Svc.Query(input)
	if err != nil {
					fmt.Println("Error getting composite index item: " , err.Error())
					fmt.Println(primary + " = :v1 AND " + sort + " = :v2")
	}

	return

}

func DeleteCompositeIndexItem(primaryValue string, sortValue string, primary string, sort string, table string) {

	deleteRequest := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
					primary: {
							S: aws.String(primaryValue),
					},
					sort: {
							S: aws.String(sortValue),
					},
			},
			TableName: aws.String(table),
	}

	_, err := Svc.DeleteItem(deleteRequest)

	if err != nil {
						fmt.Println("Got error calling DeleteItem")
						fmt.Println(err.Error())
				}


}

func DeletePrimaryItem(primaryValue string, primary string, table string, attrname string, attrvalue string) {

	deleteRequest := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
					primary: {
							S: aws.String(primaryValue),
					},
			},
			TableName: aws.String(table),
			ConditionExpression: aws.String(attrname + " = :v1"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v1": { S: aws.String(attrvalue)},
			},
	}

	_, err := Svc.DeleteItem(deleteRequest)

	if err != nil {
						fmt.Println("Got error calling DeleteItem")
						fmt.Println(err.Error())
				}


}
