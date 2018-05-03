package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  //"os"
  "github.com/google/uuid"

)

type Rank struct {
  Rid           string        `dynamodbav:"rid"`
  Title         string        `dynamodbav:"title"`
  Owner         string        `dynamodbav:"ownerid"`
  Hd            string        `dynamodbav:"hd"`
  Description   string        `dynamodbav:"description"`
  Options       []string      `dynamodbav:"options"`
  BrierScore    float64       `dynamodbav:"brierscore"`
  Concluded     bool          `dynamodbav:"concluded"`
  ConcludedTime string        `dynamodbav:"concludetime"`
}

func CreateRank (title string, description string, options []string,  hd string, owner string) (rid string){

  		ruuid := uuid.New()
  		item := Rank{
  				Rid: ruuid.String(),
          Owner: owner,
          Options: options,
          Hd: hd,
  		    Title: title,
  		    Description: description,
  		}

  		PutItem(item, "ranks-tf")

  		fmt.Println("Successfully added.")

      return ruuid.String()

}

func GetRank (rid string) (r Rank) {

  result := GetPrimaryItem(rid, "rid", "ranks-tf")

  r = Rank{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if r.Rid == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return r

}

func ListRanks(user string) (r []Rank) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", "ranks-tf")

  r = []Rank{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return r

}

func DeleteRank(rid string, owner string) {

  DeletePrimaryItem(rid, "rid", "ranks-tf", "ownerid", owner)

  fmt.Println("Deleted rank.", rid)

//  DeleteRankRanges(rid)

}

func ViewRankResults (rid string) (s []Sort) {
  //Need to do a HD check here to prevent IDOR.

    result := GetPrimaryIndexItem(rid, "rid", "rid-index", "sorts-tf")

    s = []Sort{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return s
}
