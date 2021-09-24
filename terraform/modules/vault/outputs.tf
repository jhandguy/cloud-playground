output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "root_token" {
  value       = random_password.root_token.result
  description = "Root token"
  sensitive   = true
}
