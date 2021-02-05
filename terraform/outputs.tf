output "aws_s3_endpoint" {
  value       = "http://${var.node_ip}:${module.minikube.node_ports["localstack"]}"
  description = "S3 endpoint"
}

output "aws_s3_bucket" {
  value       = module.localstack.aws_s3_buckets["s3"]
  description = "S3 bucket"
}

output "s3_host" {
  value       = var.node_ip
  description = "S3 host"
}

output "s3_port" {
  value       = module.minikube.node_ports["s3"]
  description = "S3 port"
}

output "s3_token" {
  value       = module.s3.token
  description = "S3 token"
  sensitive   = true
}

output "aws_dynamo_endpoint" {
  value       = "http://${var.node_ip}:${module.minikube.node_ports["localstack"]}"
  description = "Dynamo endpoint"
}

output "aws_dynamo_table" {
  value       = module.localstack.aws_dynamo_tables["dynamo"]
  description = "Dynamo table"
}

output "dynamo_host" {
  value       = var.node_ip
  description = "Dynamo host"
}

output "dynamo_port" {
  value       = module.minikube.node_ports["dynamo"]
  description = "Dynamo port"
}

output "dynamo_token" {
  value       = module.dynamo.token
  description = "Dynamo token"
  sensitive   = true
}