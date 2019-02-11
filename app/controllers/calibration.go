package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"
	"math"
	"strconv"
	"strings"
)

type Calibration struct {
	*revel.Controller
}

func (c Calibration) Index() revel.Result {
	return c.Render()
}

// TODO: Change these to (30, 5) for production
var sessionSize = 5
var batchSize = 2

func (c Calibration) Create() revel.Result {

	fmt.Println("Making a new session")
	sid := models.CreateCalibrationSession(c.Session["user"], sessionSize)

	c.Flash.Success(sid)

	return c.Redirect(Calibration.NextQuestions, sid)
}

func (c Calibration) NextQuestions(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)
	currentQuestionIndex := session.CurrentQuestionIndex

	// There are no more questions; redirect to the review page.
	if currentQuestionIndex >= len(session.Questions) {
		return c.Redirect(Calibration.Review, sid)
	}

	// There are some questions left, but not enough for a full batch. Return as many as there are left.
	nextBatchSize := getNextBatchSize(session)

	// Slice out the questions for this batch
	questions := session.Questions[currentQuestionIndex : currentQuestionIndex + nextBatchSize]

	return c.Render(session, questions, currentQuestionIndex)
}

func (c Calibration) SaveAnswers(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)

	result := models.GetCalibrationResult(session.ResultsId)
	fmt.Println(result)
	if len(result.Answers) == 0 {
		result.Answers = make([]models.CalibrationAnswer, len(session.Questions))
	}

	for formKey, formValue := range c.Params.Form {
		if strings.HasPrefix(formKey, "calibration-answer-") {
			questionId := strings.Replace(formKey, "calibration-answer-", "", 1)
			//question := models.GetCalibrationQuestion(questionId)

			var questionIndex int
			for questionIndex = 0; questionIndex < len(session.Questions); questionIndex++ {
				if session.Questions[questionIndex].Id == questionId {
					break
				}
			}

			// Normalize the response to be 0-100% confidence of True, for better graphing
			normalizedAnswer := true
			var err error
			var normalizedConfidence float64
			if formValue[0] == "true" {
				normalizedConfidence, err = strconv.ParseFloat(c.Params.Form.Get("calibration-confidence-" + questionId), 64)
				if err != nil {
					panic(fmt.Sprintf("Got an error parsing the confidence value, %v", err))
				}
			} else {
				normalizedConfidence, err = strconv.ParseFloat(c.Params.Form.Get("calibration-confidence-" + questionId), 64)
				normalizedConfidence = 1.0 - normalizedConfidence
				if err != nil {
					panic(fmt.Sprintf("Got an error parsing the confidence value, %v", err))
				}
			}

			normalizedOutcome := normalizedAnswer == session.Questions[questionIndex].CorrectAnswer
			//confidence, _ := strconv.ParseFloat(c.Params.Form.Get("calibration-confidence-" + questionId), 64)

			result.Answers[questionIndex] = models.CalibrationAnswer {
				Outcome:    normalizedOutcome,
				Confidence: normalizedConfidence,
			}

			fmt.Println("brier score:", brierScore(result.Answers[questionIndex]))
		}
	}

	err := models.WriteCalibrationResult(result)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	// Update the CurrentQuestionIndex of the session
	nextBatchSize := getNextBatchSize(session)
	models.UpdateCalibrationSession(session.Id, session.CurrentQuestionIndex+nextBatchSize, c.Session["user"])
	session.CurrentQuestionIndex = session.CurrentQuestionIndex+nextBatchSize

	if session.CurrentQuestionIndex < len(session.Questions) {
		return c.Redirect(Calibration.NextQuestions, sid)
	} else {
		return c.Redirect(Calibration.Review, sid)
	}
}

func getNextBatchSize(session models.CalibrationSession) int {
	var nextBatchSize int
	if session.CurrentQuestionIndex+batchSize >= len(session.Questions) {
		nextBatchSize = len(session.Questions) - session.CurrentQuestionIndex
	} else {
		nextBatchSize = batchSize
	}
	return nextBatchSize
}

func (c Calibration) Review(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)
	result := models.GetCalibrationResult(session.ResultsId)

	brierScores := make([]float64, len(result.Answers))
	fmt.Println("num answers:", len(result.Answers))
	for questionIndex := 0; questionIndex < len(result.Answers); questionIndex++ {
		fmt.Println("adding a brier score", questionIndex)
		brierScores[questionIndex] = brierScore(result.Answers[questionIndex])
	}

	fmt.Println("brier scores", brierScores)

	rawPageData := make(map[string]interface{})
	rawPageData["brierScores"] = brierScores
	rawPageData["answers"] = result.Answers

	jsonBytesPageData, _ := json.Marshal(rawPageData)
	
	pageData := string(jsonBytesPageData)

	fmt.Println(pageData)

	return c.Render(pageData)
}

func brierScore(answer models.CalibrationAnswer) float64 {
	var outcomeNum float64
	if answer.Outcome == true {
		outcomeNum = 1
	} else {
		outcomeNum = 0
	}
	score := math.Pow(outcomeNum- answer.Confidence, 2) + math.Pow((1 -outcomeNum) - (1 - answer.Confidence), 2)
	return score
}
