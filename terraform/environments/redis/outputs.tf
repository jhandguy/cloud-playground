output "postgres_url" {
  value       = module.postgresql.url
  description = "Postgres URL"
}

output "postgres_database" {
  value       = module.postgresql.database_name
  description = "Postgres database"
}

output "postgres_password" {
  value       = module.postgresql.postgres_password
  description = "Postgres password"
  sensitive   = true
}

output "mysql_url" {
  value       = module.mysql.url
  description = "MySQL URL"
}

output "mysql_database" {
  value       = module.mysql.database_name
  description = "MySQL database"
}

output "mysql_password" {
  value       = module.mysql.root_password
  description = "MySQL password"
  sensitive   = true
}

output "redis_url" {
  value       = module.redis.url
  description = "Redis URL"
}

output "redis_password" {
  value       = module.redis.redis_password
  description = "Redis password"
  sensitive   = true
}

output "sql_postgres_url" {
  value       = module.sql_postgres.url
  description = "SQL Postgres URL"
}

output "sql_postgres_token" {
  value       = random_password.sql_postgres_token.result
  description = "SQL Postgres token"
  sensitive   = true
}

output "sql_mysql_url" {
  value       = module.sql_mysql.url
  description = "SQL MySQL URL"
}

output "sql_mysql_token" {
  value       = random_password.sql_mysql_token.result
  description = "SQL MySQL token"
  sensitive   = true
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
