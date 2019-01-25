package models

import (
	"fmt"
	"math"
)

func BrierCalc(af []float64, index int) (score float64) {

	for i, v := range af {

		outcome := 0

		//If this forecast was the true outcome, mark it as such.
		if i == index {
			outcome = 1
		}

		// Square the error, add to the score.
		// Because we store percentages as whole numbers, need to convert. (99 instead of .99)
		score += math.Pow(float64(outcome)-(v*.01), 2)

	}
	fmt.Println("Current score: ", score)
	return score

}

func BrierCalcEstimate(min float64, max float64, actual float64) float64 {

	if (actual >= min) && (actual <= max) {
		return .01
	} else {
		return .81
	}

}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
