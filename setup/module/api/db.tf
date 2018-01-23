resource "aws_dynamodb_table" "tags" {
  name           = "l0-${var.name}-tags"
  read_capacity  = 50
  write_capacity = 5
  hash_key       = "EntityType"
  range_key      = "EntityID"
  tags           = "${var.tags}"

  attribute {
    name = "EntityType"
    type = "S"
  }

  attribute {
    name = "EntityID"
    type = "S"
  }
}

resource "aws_dynamodb_table" "lock" {
  name           = "l0-${var.name}-lock"
  read_capacity  = 25
  write_capacity = 10
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
