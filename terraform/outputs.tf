output "s3_endpoint" {
  value       = "http://${var.minikube_ip}:${random_integer.localstack_node_port.result}"
  description = "S3 Endpoint"
}

output "s3_bucket" {
  value       = aws_s3_bucket.s3.id
  description = "S3 Bucket"
}

output "s3_host" {
  value       = "${var.minikube_ip}:${random_integer.s3_node_port.result}"
  description = "S3 Host"
}

output "s3_token" {
  value       = random_password.s3_token.result
  description = "S3 Token"
  sensitive   = true
}

output "dynamo_endpoint" {
  value       = "http://${var.minikube_ip}:${random_integer.localstack_node_port.result}"
  description = "Dynamo Endpoint"
}

output "dynamo_table" {
  value       = aws_dynamodb_table.dynamo.id
  description = "Dynamo Table"
}

output "dynamo_host" {
  value       = "${var.minikube_ip}:${random_integer.dynamo_node_port.result}"
  description = "Dynamo Host"
}

output "dynamo_token" {
  value       = random_password.dynamo_token.result
  description = "Dynamo Token"
  sensitive   = true
}