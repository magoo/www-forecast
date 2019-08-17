package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	// Id             string `dynamodbav:"id"`
	OauthProvider  string `dynamodbav:"oauthprovider"`
	OauthID        string `dynamodbav:"oauthid"`
	FirstLoginDate string `dynamodbav:"firstlogindate"`
	Email          string `dynamodbav:"email"`
}

func GetUserByOAuth(oauthId string, oauth string) (User, bool) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"oauthid": {
				S: aws.String(oauthId),
			},
			"oauthprovider": {
				S: aws.String(oauth),
			},
		},
		TableName: aws.String(userTable),
	}

	result, err := Svc.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	user := User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	needs_creating := false
	if err != nil {
		needs_creating = true
	}
        if user.Email == "" {
                needs_creating = true
        }

	return user, needs_creating

}

func CreateUser(email string, oauthUserId string, oauth string) {
	user := User{
		Email:          email,
		OauthID:        oauthUserId,
		OauthProvider:  oauth,
		FirstLoginDate: time.Now().String(),
	}
	err := PutItem(user, userTable)
	if err != nil {
		fmt.Println(err.Error())
	}
}
