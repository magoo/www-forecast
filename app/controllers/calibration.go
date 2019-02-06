package controllers

import (
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

func (c Calibration) Create() revel.Result {

	fmt.Println("Making a new session")
	sid := models.CreateCalibrationSession()
	models.CreateEstimate("a", "a", "a", "a", "a")

	c.Flash.Success(sid)

	return c.Redirect(Calibration.NextQuestions, sid)
}

func (c Calibration) NextQuestions(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)
	questions := session.Questions
	return c.Render(session, questions)
}

func (c Calibration) SaveAnswers(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)
	//numQuestions := len(session.Questions)
	//for questionI := 0; questionI < numQuestions; questionI++ {
	//	question := session.Questions[questionI]
	//	answer := c.Params.Form.Get("calibration-answer-" + question.Id)
	//
	//	fmt.Println("Question in POST:", question.Description)
	//	if answer != "" {
	//		fmt.Println("Answer in POST:", answer)
	//		models.GetCalibrationResult(session.ResultsId)
	//	} else {
	//		fmt.Println("Answer in POST:", "NONE")
	//	}
	//}


	result := models.GetCalibrationResult(session.ResultsId)
	fmt.Println(result)
	if len(result.Answers) == 0 {
		result.Answers = make([]models.CalibrationAnswer, len(session.Questions)) // TODO: un-magicify this number
	}

	for formKey, formValue := range c.Params.Form {
		if strings.HasPrefix(formKey, "calibration-answer-") {
			questionId := strings.Replace(formKey, "calibration-answer-", "", 1)
			//question := models.GetCalibrationQuestion(questionId)

			questionIndex := -1
			for questionIndex = 0; questionIndex < len(session.Questions); questionIndex++ {
				if session.Questions[questionIndex].Id == questionId {
					break
				}
			}

			outcome := (formValue[0] == "True") == session.Questions[questionIndex].CorrectAnswer
			confidence, _ := strconv.ParseFloat(c.Params.Form.Get("calibration-confidence-" + questionId), 64)

			result.Answers[questionIndex] = models.CalibrationAnswer {
				Outcome: outcome,
				Confidence: confidence,
			}

			fmt.Println("brier score:", brierScore(result.Answers[questionIndex]))
		}
	}

	err := models.WriteCalibrationResult(result)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	//fmt.Println("Post params:", c.Params.Form.Encode())

	//return c.Redirect(Calibration.NextQuestions, sid)
	return c.Redirect(Calibration.Review, sid)
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

	return c.Render(brierScores)
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
