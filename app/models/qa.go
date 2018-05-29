package models

//I don't imagine anyone ever directly creating a struct. Lowercase?

type question struct {
  Id            string        `dynamodbav:"id"` //Uniquely identify the question
  OwnerID       string        `dynamodbav:"ownerid"` // Owner of the question, is moderator.
  Hd            string        `dynamodbav:"hd"` // Owning organization (for larger group visibility)
  Title         string        `dynamodbav:"title"`
  Description   string        `dynamodbav:"description"`
  BrierScore    float64       `dynamodbav:"brierscore"` // Any rolling Brier score we are calculating
  Concluded     bool          `dynamodbav:"concluded"` // Has this scenario shut down?
  ConcludedTime string        `dynamodbav:"concludetime"` //If so, when?
  Records       []string      `dynamodbav:"records"`    // Audit records on the scenario.
}

type answer struct {
  OwnerID       string      `dynamodbav:"ownerid"` // Owner of the question, is moderator.
  Qid           string        `dynamodbav:"id"` // The question this answers
  Hd            string        `dynamodbav:"hd"`
  Date          string        `dynamodbav:"date"`
  UserAlias     string        `dynamodbav:"useralias"` // The users fake name
}
