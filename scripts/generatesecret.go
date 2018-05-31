//Quick script to generate revel session secret
//`go run generatesecret.go`

package main

import (
  "fmt"
  "math/rand"
  "time"
)

const alphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func init() {
  rand.Seed(time.Now().UnixNano())
}

func main() {

  fmt.Println(generateSecret())
}

func generateSecret() string {
	chars := make([]byte, 64)
	for i := 0; i < 64; i++ {
		chars[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(chars)
}
