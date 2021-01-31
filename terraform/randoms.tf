resource "random_integer" "localstack_node_port" {
  min = 30000
  max = 32767
}

resource "random_pet" "s3_bucket" {}

resource "random_password" "s3_token" {
  length = 32
}

resource "random_integer" "s3_node_port" {
  min = 30000
  max = 32767
}

resource "random_pet" "dynamo_table" {}

resource "random_password" "dynamo_token" {
  length = 32
}

resource "random_integer" "dynamo_node_port" {
  min = 30000
  max = 32767
}