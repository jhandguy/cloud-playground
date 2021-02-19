output "api_key" {
  value       = random_password.gateway_api_key.result
  sensitive   = true
  description = "API key"
}

output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}