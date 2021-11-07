resource "random_password" "s3_token" {
  length  = 32
  special = false
}

resource "random_password" "dynamo_token" {
  length  = 32
  special = false
}

resource "random_password" "gateway_token" {
  length  = 32
  special = false
}

resource "random_pet" "gateway_host" {
  length    = 3
  separator = "."
}
