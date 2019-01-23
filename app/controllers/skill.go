package controllers

import (
	"fmt"
	"github.com/magoo/www-forecast/app/models"

	"github.com/revel/revel"
)

type Skill struct {
	*revel.Controller
}

func (c Skill) Index() revel.Result {

	answers := models.ListConcludedAnswers(c.Session["user"])

	var average float64 = 0

	for _, v := range answers {
		average += v.BrierScore
		fmt.Println("calc")
	}

	average = average / float64(len(answers))

	// Load this with the average Brier Score
	BrierScore := average

	// Load this with the number of concluded forecasts
	fs := len(answers)

	questions := []models.Question{}

	for _, v := range answers {
		q := models.GetQuestion(v.Id)
		questions = append(questions, q)
	}

	return c.Render(BrierScore, fs, answers, questions)
}
