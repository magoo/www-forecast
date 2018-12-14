package controllers

import (
	"github.com/revel/revel"
)

type My struct {
	*revel.Controller
}

func (c My) Index() revel.Result {

	// We need to query for answers from the logged in user.
	// Those answers each link to a question.
	// We need the concluded questions.
	// Each concluded question will have an "result index" which represents the true outcome.
	// We calculate the Brier Score using the result index from the concluded questions, on this users answers.
	// Each answer should have a brier score.
	// We display the average.

	// Load this with the average Brier Score
	bs := .2001

	// Load this with the number of concluded forecasts
	fs := 6

	return c.Render(bs, fs)
}
