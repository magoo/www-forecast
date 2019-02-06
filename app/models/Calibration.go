package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type CalibrationQuestion struct {
	Id            string   `dynamodbav:"id"`      // Uniquely identify the question
	//OwnerID       string   `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	//Date          string   `dynamodbav:"date"`
	//Hd            string   `dynamodbav:"hd"` // Owning organization (for larger group visibility)
	Description      string   `dynamodbav:"description"`
	//Concluded     bool     `dynamodbav:"concluded"`         // Has this scenario shut down?
	//ConcludedTime string   `dynamodbav:"concludetime"`      // If so, when?
	//Records       []string `dynamodbav:"records,stringset"` // Audit records on the scenario.
	//URL           string   `dynamodbav:"url"`
	CorrectAnswer bool `dynamodbav:"correctanswer"`
	Type          string   `dynamodbav:"type"`
}

type CalibrationAnswer struct {
	Outcome    bool    `dynamodbav:"outcome"` // Whether their answer was correct
	Confidence float64 `dynamodbav:"confidence"`
}

type CalibrationResult struct {
	Id        string              `dynamodbav:"id"`      // The question this answers
	//OwnerID   string              `dynamodbav:"ownerid"` // The user that did this calibration session
	//Hd        string              `dynamodbav:"hd"`
	Date      string              `dynamodbav:"date"`
	//UserAlias string              `dynamodbav:"useralias"` // The user's fake name
	//URL       string              `dynamodbav:"url"`
	Answers   []CalibrationAnswer `dynamodbav:"answers"`
}

type CalibrationSession struct {
	Id                   string                `dynamodbav:"id"`                // ID for the session
	Questions            []CalibrationQuestion `dynamodbav:"questions"`         // List of questions used for the session
	CurrentQuestionIndex int8                  `dynamodvav:"currentbatchindex"` // Index of the next question to be served from the array of questions
	ResultsId            string                `dynamodbav:"results"`
}

func GetCalibrationSession(sid string) (calibration_session CalibrationSession) {

	result := GetPrimaryItem(sid, "id", calibrationSessionTable)

	calibration_session = CalibrationSession{}

	err := dynamodbattribute.UnmarshalMap(result.Item, &calibration_session)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return calibration_session

}

func CreateCalibrationSession() (eid string) {
	euuid := uuid.New()

	// Get all calibration questions TODO: Move to another function

	session := CalibrationSession{
		Id: euuid.String(),
		// TODO: Add userID (need to figure out where to get it from)
		Questions: ListCalibrationQuestions(),
		ResultsId: CreateCalibrationResult(),
	}

	err := PutItem(session, calibrationSessionTable)

	if err != nil {
		panic(fmt.Sprintf("Error writing to db."))
	} else {
		fmt.Println("Successfully added.")
	}

	return euuid.String()
}

func CreateCalibrationResult() (eid string) {
	euuid := uuid.New()

	// Get all calibration questions TODO: Move to another function

	results := CalibrationResult{
		Id: euuid.String(),
	}

	err := PutItem(results, calibrationResultTable)

	if err != nil {
		panic(fmt.Sprintf("Error writing to db."))
	} else {
		fmt.Println("Successfully added.")
	}

	return euuid.String()
}

func GetCalibrationResult(id string) (q CalibrationResult) {

	result := GetPrimaryItem(id, "id", calibrationResultTable)

	q = CalibrationResult{}

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

func ListCalibrationQuestions() (s []CalibrationQuestion) {

	result := GetAllItems(calibrationQuestionTable)

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
