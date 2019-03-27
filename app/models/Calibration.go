package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

type CalibrationQuestion struct {
	Id            string `dynamodbav:"id"` // Uniquely identify the question
	DateUploaded  string `dynamodbav:"date"`
	Description   string `dynamodbav:"description"`
	CorrectAnswer bool   `dynamodbav:"correctanswer"`
	AnswerDetail  string `dynamodbav:"answerdetail"`
	AnswerSource  string `dynamodbav:"answersource"`
	Difficulty    string `dynamodbav:"difficulty"`
}

type CalibrationAnswer struct {
	Answer     bool    `dynamodbav:"answer"`  // The true/false answer they provided
	Outcome    bool    `dynamodbav:"outcome"` // Whether their answer was correct
	Confidence float64 `dynamodbav:"confidence"`
}

type CalibrationResult struct {
	Id      string              `dynamodbav:"id"`      // The question this answers
	OwnerID string              `dynamodbav:"ownerid"` // The user that did this calibration session
	Date    string              `dynamodbav:"date"`
	Answers []CalibrationAnswer `dynamodbav:"answers"`
}

type CalibrationSession struct {
	Id                   string                `dynamodbav:"id"`                   // ID for the session
	OwnerID              string                `dynamodbav:"ownerid"`              // The user that did this calibration session
	Questions            []CalibrationQuestion `dynamodbav:"questions"`            // List of questions used for the session
	CurrentQuestionIndex int                   `dynamodbav:"currentquestionindex"` // Index of the next question to be served from the array of questions
	ResultsId            string                `dynamodbav:"resultsid"`
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

func CreateCalibrationSession(ownerId string, numberOfQuestions int) (id string) {
	uuid := uuid.New()

	session := CalibrationSession{
		Id:                   uuid.String(),
		OwnerID:              ownerId,
		Questions:            ListCalibrationQuestions(numberOfQuestions),
		CurrentQuestionIndex: 0,
		ResultsId:            CreateCalibrationResult(ownerId),
	}

	err := PutItem(session, calibrationSessionTable)

	if err != nil {
		panic(fmt.Sprintf("Error writing to db."))
	} else {
		fmt.Println("Successfully added.")
	}

	return uuid.String()
}

func UpdateCalibrationSession(id string, currentQuestionIndex int, user string) {
	//Primary key for update query
	key := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	expressionattrvalues := map[string]*dynamodb.AttributeValue{
		":qidx": {
			N: aws.String(strconv.Itoa(currentQuestionIndex)),
		},
		":user": {
			S: aws.String(user),
		},
	}

	updateexpression := "SET currentquestionindex = :qidx"
	conditionexpression := "ownerid = :user"

	err := UpdateItem(key, updateexpression, expressionattrvalues, calibrationSessionTable, conditionexpression)
	if err != nil {
		panic(fmt.Sprintf("Failed to UpdateCalibrationSession, %v", err))
	}
	return
}

func CreateCalibrationResult(ownerId string) (id string) {
	uuid := uuid.New()
	t := time.Now()

	results := CalibrationResult{
		Id:      uuid.String(),
		OwnerID: ownerId,
		Date:    t.Format("2006-01-02"),
	}

	err := PutItem(results, calibrationResultTable)

	if err != nil {
		panic(fmt.Sprintf("Error writing to db."))
	} else {
		fmt.Println("Successfully added.")
	}

	return uuid.String()
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

func UpdateCalibrationQuestion(id string, description string, correctAnswer bool, answerDetail string, answerSource string, difficulty string) {
	//Key for the table
	key := map[string]*dynamodb.AttributeValue {
		"id": {
			S: aws.String(id),
		},
	}

	//Make our list of "expressions"
	expressionattrvalues:= map[string]*dynamodb.AttributeValue {
		//":date": {
		//
		//},
		":description": {
			S: aws.String(description),
		},
		":correctanswer": {
			S: aws.String(strconv.FormatBool(correctAnswer)),
		},
		":answerdetail": {
			S: aws.String(answerDetail),
		},
		":answersource": {
			S: aws.String(answerSource),
		},
		":difficulty": {
			S: aws.String(difficulty),
		},
	}

	fmt.Println(expressionattrvalues)

	//Case issue. Options has mixed case in other tables, should fix on production launch. See #24
	updateexpression := "SET description = :description, correctanswer = :correctanswer, answerdetail = :answerdetail, answersource = :answersource, difficulty = :difficulty"

	//Enforce moderator
	//conditionexpression := "ownerid = :user"
	conditionexpression := ""

	err := UpdateItem(key, updateexpression, expressionattrvalues, calibrationQuestionTable, conditionexpression)
	if err != nil {
		fmt.Println("Error updating calibration question:", err)
	}

	fmt.Println("Updated CalibrationQuestion.")
}


func CreateCalibrationQuestion(id string, description string, correctAnswer bool, answerDetail string, answerSource string, difficulty string) {
	t := time.Now()

	results := CalibrationQuestion{
		Id:            id,
		DateUploaded:  t.Format("2006-01-02"),
		Description:   description,
		CorrectAnswer: correctAnswer,
		AnswerDetail:  answerDetail,
		AnswerSource:  answerSource,
		Difficulty:    difficulty, // TODO: Enforce that it's one of the accepted values
	}

	err := PutItem(results, calibrationQuestionTable)

	if err != nil {
		panic(fmt.Sprintf("Error writing to db."))
	} else {
		fmt.Println("Successfully added.")
	}
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

func GetAllCalibrationQuestions() (allCalibrationQuestions []CalibrationQuestion) {
	result := GetAllItems(calibrationQuestionTable)

	//allCalibrationQuestions := []CalibrationQuestion{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &allCalibrationQuestions)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return allCalibrationQuestions
}

func ListCalibrationQuestions(numberOfQuestions int) (batchOfQuestions []CalibrationQuestion) {
	result := GetAllItems(calibrationQuestionTable)

	allCalibrationQuestions := []CalibrationQuestion{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &allCalibrationQuestions)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	rand.Seed(time.Now().UnixNano())
	shuffledIndexes := rand.Perm(len(allCalibrationQuestions))

	batchSize := numberOfQuestions;
	if batchSize > len(allCalibrationQuestions) {
		batchSize = len(allCalibrationQuestions)
	}
	batchOfQuestions = make([]CalibrationQuestion, batchSize)
	for i := 0; i < batchSize; i ++ {
		batchOfQuestions[i] = allCalibrationQuestions[shuffledIndexes[i]]
	}

	return batchOfQuestions
}

func WriteCalibrationResult(item interface{}) (err error) {
	err = PutItem(item, calibrationResultTable)

	return
}
