locals {
  helm_releases = ["localstack", "s3", "dynamo"]
  node_ports = {
    for index in range(0, length(local.helm_releases)) : local.helm_releases[index] => random_shuffle.random_node_ports.result[index]
  }
}

resource "random_shuffle" "random_node_ports" {
  result_count = length(local.helm_releases)
  input        = range(30000, 30000 + length(local.helm_releases))
}

resource "random_pet" "s3_bucket" {}

resource "random_password" "s3_token" {
  length = 32
}

resource "random_pet" "dynamo_table" {}

resource "random_password" "dynamo_token" {
  length = 32
}