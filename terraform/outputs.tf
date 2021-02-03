output "aws_s3_endpoint" {
  value       = "http://${var.node_ip}:${local.node_ports["localstack"]}"
  description = "S3 Endpoint"
}

output "aws_s3_bucket" {
  value       = aws_s3_bucket.s3.id
  description = "S3 Bucket"
}

output "s3_host" {
  value       = var.node_ip
  description = "S3 Host"
}

output "s3_port" {
  value       = local.node_ports["s3"]
  description = "S3 Port"
}

output "s3_token" {
  value       = random_password.s3_token.result
  description = "S3 Token"
  sensitive   = true
}

output "aws_dynamo_endpoint" {
  value       = "http://${var.node_ip}:${local.node_ports["localstack"]}"
  description = "Dynamo Endpoint"
}

output "aws_dynamo_table" {
  value       = aws_dynamodb_table.dynamo.id
  description = "Dynamo Table"
}

output "dynamo_host" {
  value       = var.node_ip
  description = "Dynamo Host"
}

output "dynamo_port" {
  value       = local.node_ports["dynamo"]
  description = "Dynamo Port"
}

output "dynamo_token" {
  value       = random_password.dynamo_token.result
  description = "Dynamo Token"
  sensitive   = true
}