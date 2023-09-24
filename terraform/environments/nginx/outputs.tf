output "aws_s3_endpoint" {
  value       = module.localstack.aws_s3_endpoint
  description = "S3 endpoint"
}

output "aws_s3_bucket" {
  value       = module.localstack.aws_s3_buckets["s3"]
  description = "S3 bucket"
}

output "s3_token" {
  value       = random_password.s3_token.result
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
  value       = random_password.dynamo_token.result
  description = "Dynamo token"
  sensitive   = true
}

output "dynamo_url" {
  value       = module.dynamo.url
  description = "Dynamo URL"
}

output "gateway_token" {
  value       = random_password.gateway_token.result
  description = "Gateway token"
  sensitive   = true
}

output "canary_gateway_url" {
  value       = module.gateway.urls["canary"]
  description = "Canary Gateway URL"
}

output "stable_gateway_url" {
  value       = module.gateway.urls["stable"]
  description = "Stable Gateway URL"
}

output "ingress_gateway_url" {
  value       = module.nginx.url
  description = "Gateway Ingress URL"
}

output "ingress_gateway_host" {
  value       = module.gateway.host
  description = "Gateway Ingress host"
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

output "pushgateway_url" {
  value       = module.pushgateway.url
  description = "PushGateway URL"
}

output "grafana_admin_password" {
  value       = module.prometheus.grafana_admin_password
  description = "Grafana admin password"
  sensitive   = true
}

output "nginx_url" {
  value       = module.nginx.url
  description = "NGINX Controller URL"
}

output "argorollouts_url" {
  value       = var.argorollouts_enabled ? module.argorollouts[0].url : null
  description = "NGINX Controller URL"
}
