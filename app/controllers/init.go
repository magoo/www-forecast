package controllers

import (
	"www-forecast/app/models"
	"github.com/revel/revel"
)
func init() {
	revel.OnAppStart(models.DbConnect)
}
