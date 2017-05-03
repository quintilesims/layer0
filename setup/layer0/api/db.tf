resource "aws_dynamodb_table" "tags" {
  name           = "l0-${var.name}-tags"
  read_capacity  = 100
  write_capacity = 50
  hash_key       = "TagID"

  attribute {
    name = "TagID"
    type = "N"
  }
}

resource "aws_dynamodb_table" "jobs" {
  name           = "l0-${var.name}-jobs"
  read_capacity  = 25
  write_capacity = 10
  hash_key       = "JobID"

  attribute {
    name = "JobID"
    type = "S"
  }
}
