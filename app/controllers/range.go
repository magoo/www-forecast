package controllers

import (
	"github.com/revel/revel"
	"github.com/magoo/www-forecast/app/models"
	"fmt"
)

type Range struct {
	*revel.Controller
}

func (c Range) Index() revel.Result {
		return c.Render()
}

func (c Range) Create(minimum float64, maximum float64, eid string) revel.Result {

		e := models.GetEstimate(eid)
		if (e.Concluded) {
			c.Flash.Error("This has already concluded.")

			return c.Redirect("/view/estimate/%s", eid )
		}

		models.CreateRange(c.Session["user"].(string), minimum, maximum, eid, c.Session["hd"].(string))
		//fmt.Println(options[0])
		c.Flash.Success("Range submitted.")
		fmt.Println("redirecting to: ", eid)
		return c.Redirect("/view/estimate/%s", eid )
}
