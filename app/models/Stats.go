package models

import "fmt"

func BrierCalc (af []int, index int) (score float64) {

  fmt.Println("Entering Brier Calculation with index:", index)

  for i, v := range af {
    fmt.Println("Index: ", i, "Value: ", v)
    outcome := 0

    if (i == index){
      outcome = 1
      fmt.Println("Calculating winning index!")
    }
    dec := 0.01

    score += (float64(outcome) - (float64(v) * dec)) * (float64(outcome) - (float64(v) * dec))
    fmt.Println("Current score: ", score)
  }
  fmt.Println("Final score:", score)
  fmt.Println("Outcome Array length:", len(af))
  return score / float64(len(af))

}
