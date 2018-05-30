package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/google/uuid"
  "github.com/aws/aws-sdk-go/aws"
  "time"
  "strconv"
  "errors"
)

type Scenario struct {
  Question
  Options       []string         `dynamodbav:"options"`
  Results       []int            `dynamodbav:"results"`
  ResultIndex   int              `dynamodbav:"resultindex"`

}

func CreateScenario(title string, description string, options []string, hd string, owner string) (sid string){

      t := time.Now()

  		fuuid := uuid.New()
  		item := Scenario{
        Question: Question{
  				Id: fuuid.String(),
          OwnerID: owner,
          Hd: hd,
  		    Title: title,
  		    Description: description,
          Records: []string{t.Format("2006-01-02") + ": Created.", },
        },
          Options: options,
  		}

  		PutItem(item, "scenarios-tf")

  		fmt.Println("Successfully added.")

      return fuuid.String()

}

func UpdateScenario(sid string, title string, description string, options []string, user string) {

      //Start with the key for the table
      key := map[string]*dynamodb.AttributeValue {
        "id": {
          S: aws.String(sid),
        },
      }

      //Changing this into a list of attributes.
      o, _ := dynamodbattribute.MarshalList(options)

      expressionattrvalues:= map[string]*dynamodb.AttributeValue {
        ":t": {
          S: aws.String(title),
        },
        ":d": {
          S: aws.String(description),
        },
        ":o": {
          L: o,
        },
        ":user": {
          S: aws.String(user),
        },
      }

      updateexpression := "SET title = :t, description = :d, Options = :o"
      conditionexpression := "ownerid = :user"


  		UpdateItem(key, updateexpression, expressionattrvalues, "scenarios-tf", conditionexpression)

  		fmt.Println("Updated scenario.")
}

func ViewScenario(sid string) (s Scenario) {

  // I'll need to change this to make "secret link" work.
  result := GetPrimaryItem(sid, "id", "scenarios-tf")

  s = Scenario{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

if s.Question.Id == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return s

}

func ListScenarios(user string) (s []Scenario) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "scenarios-tf")

  s = []Scenario{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return s

}

func DeleteScenario(sid string, owner string) {

  DeletePrimaryItem(sid, "id", "scenarios-tf", "ownerid", owner)

  fmt.Println("Deleted scenario.", sid)

  DeleteScenarioForecasts(sid)

}

func (s Scenario) GetAverageForecasts() (avg []int, err error) {

  sr := ViewScenarioResults(s.Question.Id)

  if len(sr) < 1 {
    return []int{}, errors.New("No results.")
  }

  avg = []int{}
	size := len(sr[0].Forecasts)

	for i := 0; i < size; i++ {
		sum := 0
			for _, v := range sr {
					sum += v.Forecasts[i]
					//fmt.Println("Adding forecast: ", v.Forecasts[i])
			}
			//fmt.Println("Adding average to array: ", sum / len(sr))
		avg = append(avg, sum / len(sr))
	}

  return avg, nil

}

func (s Scenario) AddRecord(user string) {

  results, err := s.GetAverageForecasts()

  if err != nil {
    return
  }

  // func UpdateItem(key map[string]*dynamodb.AttributeValue, updateexpression string, expressionattrvalues map[string]*dynamodb.AttributeValue, table string, conditionexpression string ) (err error) {
  //Primary key for update query
  key := map[string]*dynamodb.AttributeValue {
    "id": {
      S: aws.String(s.Question.Id),
    },
  }

  t := time.Now()
  record := t.Format("2006-01-02") + ": Results recorded. "

   x := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

  for i, v := range results {
    record += "" + string(x[i%25]) + ". " + " " + strconv.Itoa(v) + "% "
  }


  item := map[string]*dynamodb.AttributeValue {
    ":r": {
        SS: []*string{
          aws.String(record),
          },
        },
    ":user": {
      S: aws.String(user),
    },
  }

  //av, err := dynamodbattribute.MarshalMap(item)

  UpdateItem(key, "ADD records :r", item, "scenarios-tf", "ownerid = :user")



}
