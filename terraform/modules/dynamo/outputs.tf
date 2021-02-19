output "token" {
  value       = random_password.dynamo_token.result
  sensitive   = true
  description = "Token"
}

output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}