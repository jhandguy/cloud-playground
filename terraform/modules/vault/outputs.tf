output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.vault.namespace}.${helm_release.vault.name}:8200"
  description = "Cluster URL"
}

output "root_token" {
  value       = random_password.root_token.result
  description = "Root token"
  sensitive   = true
}
