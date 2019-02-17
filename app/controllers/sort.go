package controllers

import (
	"github.com/revel/revel"
	"github.com/magoo/www-forecast/app/models"
	"fmt"
)

type Sort struct {
	*revel.Controller
}

func (c Sort) Index() revel.Result {
		return c.Render()
}

func (c Sort) Create(options []int, rid string) revel.Result {

		r := models.GetRank(rid)
		if (r.Concluded) {
			c.Flash.Error("This has already concluded.")

			return c.Redirect("/view/rank/%s", rid )
		}

		models.CreateSort(c.Session["user"], options, rid, c.Session["hd"])
		//fmt.Println(options[0])
		c.Flash.Success("Sort submitted.")
		fmt.Println("redirecting to: ", rid)
		return c.Redirect("/view/rank/%s", rid )
}
