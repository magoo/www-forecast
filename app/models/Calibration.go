package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type CalibrationQuestion struct {
	Id            string   `dynamodbav:"id"`      //Uniquely identify the question
	OwnerID       string   `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	Date          string   `dynamodbav:"date"`
	Hd            string   `dynamodbav:"hd"` // Owning organization (for larger group visibility)
	Question      string   `dynamodbav:"description"`
	Concluded     bool     `dynamodbav:"concluded"`         // Has this scenario shut down?
	ConcludedTime string   `dynamodbav:"concludetime"`      //If so, when?
	Records       []string `dynamodbav:"records,stringset"` // Audit records on the scenario.
	URL           string   `dynamodbav:"url"`
	Type          string   `dynamodbav:"type"`
}

type CalibrationAnswer struct {
	Outcome    bool    `dynamodbav:"outcome"`
	Confidence float64 `dynamodbav:"confidence"`
}

type CalibrationResults struct {
	Id        string              `dynamodbav:"id"`      // The question this answers
	OwnerID   string              `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	Hd        string              `dynamodbav:"hd"`
	Date      string              `dynamodbav:"date"`
	UserAlias string              `dynamodbav:"useralias"` // The users fake name
	URL       string              `dynamodbav:"url"`
	Answers   []CalibrationAnswer `dynamodbav:"answers"`
}

func GetCalibrationQuestion(id string) (q CalibrationQuestion) {

	result := GetPrimaryItem(id, "id", calibrationQuestionTable)

	q = CalibrationQuestion{}

	err := dynamodbattribute.UnmarshalMap(result.Item, &q)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if q.Id == "" {
		fmt.Println("Could not find that question.")
		return
	}

	return q
}

func ListCalibrationQuestions(user string) (s []CalibrationQuestion) {

	result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", calibrationQuestionTable)

	s = []CalibrationQuestion{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return s

}

func WriteCalibrationResult(item interface{}) (err error) {

	err = PutItem(item, calibrationResultTable)

	return

}
