package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
)

type View struct {
	*revel.Controller
}

func (c View) Index(fid string) revel.Result {
		f := models.ViewForecast(fid)
		return c.Render(f)
}
