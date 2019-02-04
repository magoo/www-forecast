package controllers

import (
	"fmt"
	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel"
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
	models.CreateEstimate("a","a","a", "a","a")

	//c.Flash.Out["session_id"] = sid
	c.Flash.Success(sid)

	return c.Redirect(Calibration.NextQuestions, sid)
}

func (c Calibration) NextQuestions(sid string) revel.Result {
	session := models.GetCalibrationSession(sid)
	questions := session.Questions
	return c.Render(session, questions)
}

func (c Calibration) SaveAnswers(sid string) revel.Result {
	fmt.Println("Post params:", c.Params.Form.Encode())

	return c.Redirect(Calibration.NextQuestions, sid)
}
