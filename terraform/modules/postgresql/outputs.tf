output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.postgresql.name}.${helm_release.postgresql.namespace}.svc.cluster.local:5432"
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

output "postgres_password" {
  value       = random_password.postgres_password.result
  description = "Postgres password"
}
