package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
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

func UpdateRank (rid string, title string, description string, options []string, user string) {

  //Key for the table
  key := map[string]*dynamodb.AttributeValue {
    "rid": {
      S: aws.String(rid),
    },
  }

  //Changing this into a list of attributes.
  o, _ := dynamodbattribute.MarshalList(options)

  //Make our list of "expressions"
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

  //Case issue. Options has mixed case in other tables, should fix on production launch. See #24
  updateexpression := "SET title = :t, description = :d, options = :o"

  //Enforce moderator
  conditionexpression := "ownerid = :user"

  UpdateItem(key, updateexpression, expressionattrvalues, "ranks-tf", conditionexpression)

  fmt.Println("Updated rank.")


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
