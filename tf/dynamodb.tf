provider "aws" {
  profile             = "magoo"
  region              = "us-west-2"
}

resource "aws_dynamodb_table" "scenarios" {
  name           = "scenarios-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "sid"

  attribute {
     name = "sid"
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

resource "aws_dynamodb_table" "forecasts" {
  name           = "forecasts-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "sid"
  range_key       = "ownerid"

  attribute {
     name = "sid"
     type = "S"
  }
  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "sid-index"
    hash_key           = "sid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}

resource "aws_dynamodb_table" "estimates" {
  name           = "estimates-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "eid"

  attribute {
     name = "eid"
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

resource "aws_dynamodb_table" "ranges" {
  name           = "ranges-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "eid"
  range_key       = "ownerid"

  attribute {
     name = "eid"
     type = "S"
  }
  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "eid-index"
    hash_key           = "eid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}

resource "aws_dynamodb_table" "ranks" {
  name           = "ranks-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "rid"

  attribute {
     name = "rid"
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

resource "aws_dynamodb_table" "sorts" {
  name           = "sorts-tf"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "rid"
  range_key       = "ownerid"

  attribute {
     name = "rid"
     type = "S"
  }
  attribute {
     name = "ownerid"
     type = "S"
  }

  global_secondary_index {
    name               = "rid-index"
    hash_key           = "rid"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}
