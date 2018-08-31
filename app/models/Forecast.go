package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/docker/docker/pkg/namesgenerator"
	//"os"
	"time"
)

type Forecast struct {
	Answer
	Forecasts []float64 `dynamodbav:"forecasts"`
}

func CreateForecast(u string, f []float64, sid string, hd string) {
	//Must do a permission check in the future to prevent crossover forecasts. Tock day.
	//Must do a check to make sure the array of values is equal to the array of options in the sid.

	t := time.Now()

	item := Forecast{
		Answer: Answer{
			Hd:        hd,
			Id:        sid,
			Date:      t.String(),
			OwnerID:   u,
			UserAlias: namesgenerator.GetRandomName(0),
		},
		Forecasts: f,
	}

	err := PutItem(item, "answers-tf")

	if err != nil {
		fmt.Println("Error writing to db.")
	} else {
		fmt.Println("Successfully added.")
	}

}

func ViewScenarioResults(sid string) (c []Forecast) {
	//Need to do a HD check here to prevent IDOR.

	result := GetPrimaryIndexItem(sid, "id", "id-index", "answers-tf")

	c = []Forecast{}

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &c)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return c
}

func ViewUserScenarioResults(uid string, sid string) (userForecast Forecast) {
	results := ViewScenarioResults(sid)
	for _, result := range results {
		if result.Answer.OwnerID == uid {
			userForecast = result
		}
	}

	return userForecast
}

func DeleteScenarioForecasts(sid string) {

	fs := ViewScenarioResults(sid)

	for _, v := range fs {
		fmt.Println("Deleting: ", v.Answer.Id, v.Answer.OwnerID)
		DeleteCompositeIndexItem(v.Answer.Id, v.Answer.OwnerID, "sid", "user", "answers-tf")
	}

	fmt.Println("Deleted forecasts associated with scenario.")

}
