package models

import (
  "fmt"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  //"os"
  "github.com/google/uuid"
  "time"
  "strconv"
  "sort"
  "errors"

)

type Rank struct {
  Question
  Options       []string      `dynamodbav:"options"`
}

func CreateRank (title string, description string, options []string,  hd string, owner string) (rid string){

      t := time.Now()

  		ruuid := uuid.New()
  		item := Rank{
          Question: Question{
            Id: ruuid.String(),
            OwnerID: owner,
            Date: t.Format("2006-01-02"),
            Hd: hd,
            Title: title,
            Description: description,
            Records: []string{t.Format("2006-01-02") + ": Created.", },
            URL: "rank/" + ruuid.String(),
            Type: "Rank",
          },
  				Options: options,
  		}

  		err := PutItem(item, questionTable)

      if err != nil {
        fmt.Println("Error writing to db.")
      } else {
        fmt.Println("Successfully added.")
      }

      return ruuid.String()

}

func (r Rank) GetURL() (url string) {

  return "/view/rank/" + r.Id

}

func UpdateRank (rid string, title string, description string, options []string, user string) {

  //Key for the table
  key := map[string]*dynamodb.AttributeValue {
    "id": {
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

  UpdateItem(key, updateexpression, expressionattrvalues, questionTable, conditionexpression)

  fmt.Println("Updated rank.")


}

func GetRank (rid string) (r Rank) {

  result := GetPrimaryItem(rid, "id", questionTable)

  r = Rank{}

  err := dynamodbattribute.UnmarshalMap(result.Item, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  if r.Question.Id == "" {
      fmt.Println("Could not find that scenario.")
      return
  }

  return r

}

func ListRanks(user string) (r []Rank) {

  result := GetPrimaryIndexItem(user, "ownerid", "ownerid-index", questionTable)

  r = []Rank{}

  err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &r)

  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return r

}

func DeleteRank(rid string, owner string) {

  DeletePrimaryItem(rid, "id", questionTable, "ownerid", owner)

  fmt.Println("Deleted rank.", rid)

  DeleteRankSorts(rid)

}

func DeleteRankSorts(rid string) {

    rr := ViewRankResults(rid)

    for _, v  := range rr {
      fmt.Println("Deleting: ", v.Answer.Id, v.Answer.OwnerID)
      DeleteCompositeIndexItem(v.Answer.Id, v.Answer.OwnerID, "id", "ownerid", "answers-tf")
    }


    fmt.Println("Deleted sorts associated with rank.")

}

func ViewRankResults (rid string) (s []Sort) {
  //Need to do a HD check here to prevent IDOR.

    result := GetPrimaryItem(rid, "id", "answers-tf")

    s = []Sort{}

    err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &s)

    if err != nil {
      panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    return s
}

func (r Rank) AddRecord(user string) (err error){

  er := ViewRankResults(r.Question.Id)
  if len(er) < 1{
    return errors.New("Trying to count position winner when there are no sorts.")
  }

  pw := GetPositionalWinner(er)

  t := time.Now()

  record := t.Format("2006-01-02") + ": Results recorded. "

  x := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

  for _, v := range pw {
    record += "" + string(x[v.Index%25]) + ": " + " " + strconv.Itoa(v.Votes) + " "
  }

  //Primary key for update query
  key := map[string]*dynamodb.AttributeValue {
    "id": {
      S: aws.String(r.Question.Id),
    },
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

  UpdateItem(key, "ADD records :r", item, questionTable, "ownerid = :user")

  return err

}

func GetPositionalWinner(rr []Sort) (vs Votes){

	vs = make(Votes, len(rr[0].Options))

	total := len(rr[0].Options)

	//First loop. 'v' is a Sort.
	for _, v := range rr {

		//Second loop. Each "o" is a preference, top to bottom.
		for i, o := range v.Options {
			vs[o].Votes += total - i - 1
			vs[o].Index = o
		}
	}

	sort.Sort(sort.Reverse(vs))

	return vs
}
