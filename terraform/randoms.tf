resource "random_password" "s3_token" {
  length  = 32
  special = false
}

resource "random_password" "dynamo_token" {
  length  = 32
  special = false
}

resource "random_password" "gateway_api_key" {
  length  = 32
  special = false
}