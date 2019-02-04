package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type CalibrationQuestion struct {
	Id            string   `dynamodbav:"id"`      // Uniquely identify the question
	OwnerID       string   `dynamodbav:"ownerid"` // Owner of the question, is moderator.
	//Date          string   `dynamodbav:"date"`
	//Hd            string   `dynamodbav:"hd"` // Owning organization (for larger group visibility)
	Question      string   `dynamodbav:"description"`
	//Concluded     bool     `dynamodbav:"concluded"`         // Has this scenario shut down?
	//ConcludedTime string   `dynamodbav:"concludetime"`      // If so, when?
	//Records       []string `dynamodbav:"records,stringset"` // Audit records on the scenario.
	//URL           string   `dynamodbav:"url"`
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
	UserAlias string              `dynamodbav:"useralias"` // The user's fake name
	URL       string              `dynamodbav:"url"`
	Answers   []CalibrationAnswer `dynamodbav:"answers"`
}

type CalibrationSession struct {
	Id                   string                `dynamodbav:"id"`                // ID for the session
	Questions            []CalibrationQuestion `dynamodbav:"questions"`         // List of questions used for the session
	CurrentQuestionIndex int8                  `dynamodvav:"currentbatchindex"` // Index of the next question to be served from the array of questions
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

	//t := time.Now()
	//
	//euuid := uuid.New()
	//item := Estimate{
	//	Question: Question{
	//		Id:          euuid.String(),
	//		Date:        t.Format("2006-01-02"),
	//		OwnerID:     owner,
	//		Hd:          hd,
	//		Title:       title,
	//		Description: description,
	//		Records:     []string{t.Format("2006-01-02") + ": Created."},
	//		URL:         "estimate/" + euuid.String(),
	//		Type:        "Estimate",
	//	},
	//	Unit: unit,
	//}
	//
	//err := PutItem(item, questionTable)
	//
	//if err != nil {
	//	fmt.Println("Error writing to db.")
	//} else {
	//	fmt.Println("Successfully added.")
	//}
	//
	//return euuid.String()

	euuid := uuid.New()

	// Get all calibration questions TODO: Move to another function

	session := CalibrationSession{
		Id: euuid.String(),
		Questions: ListCalibrationQuestions(),
	}

	err := PutItem(session, calibrationSessionTable)

	if err != nil {
		fmt.Println("Error writing to db.")
	} else {
		fmt.Println("Successfully added.")
	}

	return euuid.String()
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
