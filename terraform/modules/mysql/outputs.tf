output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.mysql.name}.${helm_release.mysql.namespace}.svc.cluster.local:3306"
  description = "Cluster URL"
}

output "database_name" {
  value       = var.database_name
  description = "Database name"
}

output "user_name" {
  value       = var.user_name
  description = "User name"
}

output "user_password" {
  value       = random_password.user_password.result
  description = "User password"
}

output "root_password" {
  value       = random_password.root_password.result
  description = "Root password"
}
