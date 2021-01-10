resource "aws_dynamodb_table" "table" {
  depends_on = [kubernetes_ingress.localstack]

  name         = random_pet.table.id
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }
}