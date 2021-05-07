output "api_key" {
  value       = random_password.gateway_api_key.result
  sensitive   = true
  description = "API key"
}

output "urls" {
  value = {
    for name, node_port in var.node_ports : name => "${var.node_ip}:${node_port}"
  }
  description = "URLs"
}