output "aws_s3_endpoint" {
  value       = module.localstack.aws_dynamo_endpoint
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

output "prod_gateway_url" {
  value       = module.gateway.urls["prod"]
  description = "Prod Gateway URL"
}

output "ingress_gateway_url" {
  value       = module.consul.ingress_gateway_urls["gateway"]
  description = "Gateway Ingress URL"
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

output "consul_url" {
  value       = module.consul.url
  description = "Consul URL"
}

output "vault_url" {
  value       = module.vault.url
  description = "Vault URL"
}

output "vault_root_token" {
  value       = module.vault.root_token
  description = "Vault root token"
  sensitive   = true
}
