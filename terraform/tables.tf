resource "aws_dynamodb_table" "dynamo" {
  depends_on = [helm_release.localstack]

  name         = random_pet.dynamo_table.id
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }
}