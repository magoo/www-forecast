package models

import (

  "math"
)

func BrierCalc (af []float64, index int) (score float64) {

  //fmt.Println("Entering Brier Calculation with index:", index)

  for i, v := range af {
    //fmt.Println("Index: ", i, "Value: ", v)
    outcome := 0

    if (i == index){
      outcome = 1
      //fmt.Println("Calculating winning index!")
    }
    dec := 0.01

    score += (float64(outcome) - (v * dec)) * (float64(outcome) - (v * dec))
    //fmt.Println("Current score: ", score)
  }

  return score / float64(len(af))

}

func BrierCalcEstimate (min float64, max float64, actual float64) float64 {

  if ((actual >= min) && (actual <= max)){
    return .01
  } else {
    return .81
  }

}

func RoundPlus(f float64, places int) (float64) {
    shift := math.Pow(10, float64(places))
    return Round(f * shift) / shift;
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
