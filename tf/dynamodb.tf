provider "aws" {
  profile             = "e6e-dev"
  region              = "us-west-1"
}

resource "aws_dynamodb_table" "questions" {
  name           = "questions-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "id"

  attribute {
     name = "id"
     type = "S"
  }

  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "ownerid-index"
    hash_key           = "ownerid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }
}

resource "aws_dynamodb_table" "answers" {
  name           = "answers-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "id"
  range_key       = "ownerid"

  attribute {
     name = "id"
     type = "S"
  }
  attribute {
     name = "ownerid"
     type = "S"
  }
  attribute {
     name = "qid"
     type = "S"
  }

  global_secondary_index {
    name               = "qid-index"
    hash_key           = "qid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}
