package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Question struct {
	Id            string   `dynamodbav:"id"`      //Uniquely identify the question
	OwnerID       string   `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	Date          string   `dynamodbav:"date"`
	Hd            string   `dynamodbav:"hd"` // Owning organization (for larger group visibility)
	Title         string   `dynamodbav:"title"`
	Description   string   `dynamodbav:"description"`
	BrierScore    float64  `dynamodbav:"brierscore"`        // Any rolling Brier score we are calculating
	Concluded     bool     `dynamodbav:"concluded"`         // Has this scenario shut down?
	ConcludedTime string   `dynamodbav:"concludetime"`      //If so, when?
	Records       []string `dynamodbav:"records,stringset"` // Audit records on the scenario.
	URL           string   `dynamodbav:"url"`
	Type          string   `dynamodbav:"type"`
}

type Answer struct {
	Id          string  `dynamodbav:"id"`      // The question this answers
	OwnerID     string  `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	Hd          string  `dynamodbav:"hd"`
	Date        string  `dynamodbav:"date"`
	UserAlias   string  `dynamodbav:"useralias"` // The users fake name
	URL         string  `dynamodbav:"url"`
	Title       string  `dynamodbav:"title"`
	Description string  `dynamodbav:"description"`
	BrierScore  float64 `dynamodbav:"brierscore"`
	Concluded   bool    `dynamodbav:"concluded"`
}

func GetQuestion(id string) (q Question) {

	result := GetPrimaryItem(id, "id", questionTable)

	q = Question{}

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

func ListQuestions(user string) (s []Question) {

	result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", questionTable)

	s = []Question{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return s

}

func ListAnswers(user string) (s []Answer) {
	result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", answerTable)

	s = []Answer{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return s

}

func ListConcludedAnswers(user string) (cs []Answer) {

	s := ListAnswers(user)

	//Concluded Questions
	cs = []Answer{}

	for _, v := range s {
		if v.Concluded {
			cs = append(cs, v)
		}
	}

	return cs

}

func (q Question) WriteRecord(message string, user string) (err error) {

	if err != nil {
		return err
	}

	// func UpdateItem(key map[string]*dynamodb.AttributeValue, updateexpression string, expressionattrvalues map[string]*dynamodb.AttributeValue, table string, conditionexpression string ) (err error) {
	//Primary key for update query
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(q.Id),
		},
	}

	t := time.Now()
	record := t.Format("2006-01-02 3:04:05PM") + ": " + message

	item := map[string]*dynamodb.AttributeValue{
		":r": {
			SS: []*string{
				aws.String(record),
			},
		},
		":user": {
			S: aws.String(user),
		},
	}

	//av, err := dynamodbattribute.MarshalMap(item)

	UpdateItem(key, "ADD records :r", item, questionTable, "ownerid = :user")

	return err

}

func ViewQuestionAnswers(id string) (as []Answer) {
	//Need to do a HD check here to prevent IDOR.

	result := GetPrimaryIndexItem(id, "id", "id-index", answerTable)

	as = []Answer{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &as)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return as
}

func DeleteQuestionAnswers(id string) (err error) {

	qas := ViewQuestionAnswers(id)

	for _, v := range qas {
		fmt.Println("Deleting: ", v.Id, v.OwnerID)
		DeleteCompositeIndexItem(v.Id, v.OwnerID, "id", "ownerid", answerTable)
	}

	fmt.Println("Deleted answers associated with question.")

	return nil

}

func WriteQuestion(item interface{}) (err error) {

	err = PutItem(item, questionTable)

	return

}

func WriteAnswer(item interface{}) (err error) {

	err = PutItem(item, answerTable)

	return
}
