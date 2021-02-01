output "aws_s3_endpoint" {
  value       = "http://${var.node_ip}:${random_integer.localstack_node_port.result}"
  description = "S3 Endpoint"
}

output "aws_s3_bucket" {
  value       = aws_s3_bucket.s3.id
  description = "S3 Bucket"
}

output "s3_host" {
  value       = "${var.node_ip}:${random_integer.s3_node_port.result}"
  description = "S3 Host"
}

output "s3_token" {
  value       = random_password.s3_token.result
  description = "S3 Token"
  sensitive   = true
}

output "aws_dynamo_endpoint" {
  value       = "http://${var.node_ip}:${random_integer.localstack_node_port.result}"
  description = "Dynamo Endpoint"
}

output "aws_dynamo_table" {
  value       = aws_dynamodb_table.dynamo.id
  description = "Dynamo Table"
}

output "dynamo_host" {
  value       = "${var.node_ip}:${random_integer.dynamo_node_port.result}"
  description = "Dynamo Host"
}

output "dynamo_token" {
  value       = random_password.dynamo_token.result
  description = "Dynamo Token"
  sensitive   = true
}