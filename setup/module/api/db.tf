resource "aws_dynamodb_table" "tags" {
  name           = "l0-${var.name}-tags"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "EntityType"
  range_key      = "EntityID"
  tags           = var.tags

  attribute {
    name = "EntityType"
    type = "S"
  }

  attribute {
    name = "EntityID"
    type = "S"
  }
  
  ttl {
    attribute_name = "TimeToExist"
    enabled        = true
  }
}

resource "aws_dynamodb_table" "jobs" {
  name           = "l0-${var.name}-jobs"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "JobID"

  attribute {
    name = "JobID"
    type = "S"
  }
  ttl {
    attribute_name = "TimeToExist"
    enabled        = true
  }
}