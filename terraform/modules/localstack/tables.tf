resource "aws_dynamodb_table" "dynamo" {
  depends_on = [helm_release.localstack]
  for_each = {
    for index in range(length(var.aws_dynamo_tables)) : var.aws_dynamo_tables[index] => random_id.tables[index].hex
  }

  name         = each.value
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }
}
