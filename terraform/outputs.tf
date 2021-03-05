output "aws_s3_endpoint" {
  value       = module.localstack.aws_dynamo_endpoint
  description = "S3 endpoint"
}

output "aws_s3_bucket" {
  value       = module.localstack.aws_s3_buckets["s3"]
  description = "S3 bucket"
}

output "s3_token" {
  value       = module.s3.token
  description = "S3 token"
  sensitive   = true
}

output "s3_url" {
  value       = module.s3.url
  description = "S3 URL"
}

output "aws_dynamo_endpoint" {
  value       = module.localstack.aws_dynamo_endpoint
  description = "Dynamo endpoint"
}

output "aws_dynamo_table" {
  value       = module.localstack.aws_dynamo_tables["dynamo"]
  description = "Dynamo table"
}

output "dynamo_token" {
  value       = module.dynamo.token
  description = "Dynamo token"
  sensitive   = true
}

output "dynamo_url" {
  value       = module.dynamo.url
  description = "Dynamo URL"
}

output "gateway_api_key" {
  value       = module.gateway.api_key
  description = "Gateway API key"
  sensitive   = true
}

output "gateway_url" {
  value       = module.gateway.url
  description = "Gateway URL"
}

output "prometheus_url" {
  value       = module.prometheus.prometheus_url
  description = "Prometheus URL"
}

output "alertmanager_url" {
  value       = module.prometheus.alertmanager_url
  description = "AlertManager URL"
}

output "grafana_url" {
  value       = module.prometheus.grafana_url
  description = "Grafana URL"
}

output "grafana_admin_password" {
  value       = module.prometheus.grafana_admin_password
  description = "Grafana admin password"
  sensitive   = true
}