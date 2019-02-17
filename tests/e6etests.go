package tests

import (
	"math"

	"github.com/magoo/www-forecast/app/models"
	"github.com/revel/revel/testing"
)

type E6etests struct {
	testing.TestSuite
}

func (t *E6etests) TestBrierScore() {

	// Whole number percentages (90, not .9) are stored in the DB.

	af := []float64{0, 100}
	index := 1

	bs := models.BrierCalc(af, index)

	//This brier score must equal zero
	println("Testing 2 forecasts.")
	t.Assert(bs == (math.Pow(0-af[0]*.01, 2) + math.Pow(1-af[1]*.01, 2)))

	af = []float64{0, 10, 50, 40}
	index = 0

	bs = models.BrierCalc(af, index)

	//This brier score must equal 1.42
	println("Testing 4 forecasts. (Totally certain and wrong)")
	t.Assert(bs == (math.Pow(1-af[0]*.01, 2) + math.Pow(0-af[1]*.01, 2) + math.Pow(0-af[2]*.01, 2) + math.Pow(0-af[3]*.01, 2)))

	af = []float64{0, 10, 50, 40}
	index = 2

	bs = models.BrierCalc(af, index)

	//This brier score must equal .42 with a long trail
	println("Tested 4 forecasts. (Sorta certain)")
	t.Assert(bs == (math.Pow(0-af[0]*.01, 2) + math.Pow(0-af[1]*.01, 2) + math.Pow(1-af[2]*.01, 2) + math.Pow(0-af[3]*.01, 2)))

}
