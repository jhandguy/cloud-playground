output "postgresql_url" {
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

output "sql_mysql_url" {
  value       = module.sql_mysql.url
  description = "SQL MySQL URL"
}
