package models

import (
		"github.com/aws/aws-sdk-go/aws"
	  "github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/aws/credentials"
	  "github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess, _ = session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "magoo"),
	},
)

var dbname = "scenarios"

var Svc = dynamodb.New(sess)

func DbConnect(){





}
