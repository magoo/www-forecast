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

// The number of questions to use for a whole calibration training session
var sessionSize = 30
// The number of questions to show at once before moving on to the next
var batchSize = 5

func (c Calibration) Create() revel.Result {
	fmt.Println("Making a new session")
	sid := models.CreateCalibrationSession(c.Session["user"], sessionSize)

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
	questions := session.Questions[currentQuestionIndex : currentQuestionIndex+nextBatchSize]

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

			var questionIndex int
			for questionIndex = 0; questionIndex < len(session.Questions); questionIndex++ {
				if session.Questions[questionIndex].Id == questionId {
					break
				}
			}

			// Normalize the response to be 0-100% confidence of True, for better graphing
			originalAnswer := formValue[0] == "true"
			var err error
			originalConfidence, err := strconv.ParseFloat(c.Params.Form.Get("calibration-confidence-"+questionId), 64)
			if err != nil {
				panic(fmt.Sprintf("Got an error parsing the confidence value, %v", err))
			}

			originalOutcome := originalAnswer == session.Questions[questionIndex].CorrectAnswer

			result.Answers[questionIndex] = models.CalibrationAnswer{
				Answer:     originalAnswer,
				Outcome:    originalOutcome,
				Confidence: originalConfidence,
			}
		}
	}

	err := models.WriteCalibrationResult(result)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	// Update the CurrentQuestionIndex of the session
	nextBatchSize := getNextBatchSize(session)
	models.UpdateCalibrationSession(session.Id, session.CurrentQuestionIndex+nextBatchSize, c.Session["user"])
	session.CurrentQuestionIndex = session.CurrentQuestionIndex + nextBatchSize

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

	numQuestions := len(result.Answers)

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

	// Calculate mean confidence
	var totalConfidence float64
	for i := 0; i < len(result.Answers); i++ {
		totalConfidence += result.Answers[i].Confidence
	}
	meanConfidence := totalConfidence / float64(len(result.Answers))
	percentMeanConfidence := meanConfidence * 100

	// Calculate percent correct
	var totalCorrect int
	for i := 0; i < len(result.Answers); i++ {
		if result.Answers[i].Outcome {
			totalCorrect += 1
		}
	}
	fractionCorrect := float64(totalCorrect) / float64(len(result.Answers))
	percentFractionCorrect := fractionCorrect * 100

	// Calculate mean confidence on correct answers
	var totalConfidenceOnCorrect float64
	for i := 0; i < len(result.Answers); i++ {
		if result.Answers[i].Outcome {
			totalConfidenceOnCorrect += result.Answers[i].Confidence
		}
	}
	meanConfidenceOnCorrect := totalConfidenceOnCorrect / float64(totalCorrect)
	percentMeanConfidenceOnCorrect := meanConfidenceOnCorrect * 100

	// Calculate mean confidence on incorrect answers
	var totalConfidenceOnIncorrect float64
	for i := 0; i < len(result.Answers); i++ {
		if !result.Answers[i].Outcome {
			totalConfidenceOnIncorrect += result.Answers[i].Confidence
		}
	}
	meanConfidenceOnIncorrect := totalConfidenceOnIncorrect / float64(totalCorrect)
	percentMeanConfidenceOnIncorrect := meanConfidenceOnIncorrect * 100

	// Calculate total low, medium, and high confidence answers
	var totalLowAnswered int
	var totalMedAnswered int
	var totalHighAnswered int
	var totalLowCorrect int
	var totalMedCorrect int
	var totalHighCorrect int
	for i := 0; i < len(result.Answers); i++ {
		answer := result.Answers[i]
		if answer.Confidence == 0.5 || answer.Confidence == 0.6 {
			totalLowAnswered ++
			if answer.Outcome {
				totalLowCorrect ++
			}
		}
		if answer.Confidence == 0.7 || answer.Confidence == 0.8 {
			totalMedAnswered ++
			if answer.Outcome {
				totalMedCorrect ++
			}
		}
		if answer.Confidence == 0.9 || answer.Confidence == 1.0 {
			totalHighAnswered ++
			if answer.Outcome {
				totalHighCorrect ++
			}
		}
	}
	fractionLowCorrect := float64(totalLowCorrect) / float64(totalLowAnswered)
	if math.IsNaN(fractionLowCorrect) {
		fractionLowCorrect = 0
	}
	fractionMedCorrect := float64(totalMedCorrect) / float64(totalMedAnswered)
	if math.IsNaN(fractionMedCorrect) {
		fractionMedCorrect = 0
	}
	fractionHighCorrect := float64(totalHighCorrect) / float64(totalHighAnswered)
	if math.IsNaN(fractionHighCorrect) {
		fractionHighCorrect = 0
	}

	return c.Render(
		pageData,
		percentMeanConfidence,
		percentFractionCorrect,
		percentMeanConfidenceOnCorrect,
		percentMeanConfidenceOnIncorrect,
		numQuestions,
		totalCorrect,
		totalLowAnswered,
		totalMedAnswered,
		totalHighAnswered,
		totalLowCorrect,
		totalMedCorrect,
		totalHighCorrect,
		fractionLowCorrect,
		fractionMedCorrect,
		fractionHighCorrect,
	)
}

func brierScore(answer models.CalibrationAnswer) float64 {
	var outcomeNum float64
	if answer.Outcome == true {
		outcomeNum = 1
	} else {
		outcomeNum = 0
	}
	score := math.Pow(outcomeNum-answer.Confidence, 2) + math.Pow((1-outcomeNum)-(1-answer.Confidence), 2)
	return score
}
