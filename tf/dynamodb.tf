provider "aws" {
  profile             = "e6e-dev"
  region              = "us-west-1"
}

resource "aws_dynamodb_table" "scenarios" {
  name           = "scenarios-tf"
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

resource "aws_dynamodb_table" "forecasts" {
  name           = "forecasts-tf"
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

  global_secondary_index {
    name               = "id-index"
    hash_key           = "id"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}

resource "aws_dynamodb_table" "estimates" {
  name           = "estimates-tf"
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

resource "aws_dynamodb_table" "ranges" {
  name           = "ranges-tf"
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

  global_secondary_index {
    name               = "id-index"
    hash_key           = "id"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}

resource "aws_dynamodb_table" "ranks" {
  name           = "ranks-tf"
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

resource "aws_dynamodb_table" "sorts" {
  name           = "sorts-tf"
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

  global_secondary_index {
    name               = "id-index"
    hash_key           = "id"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "ALL"
  }

}
