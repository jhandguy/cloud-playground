module "kind" {
  source = "../../modules/kind"

  cluster_name = "redis"
  node_ports = [
    "postgresql",
    "mysql",
    "redis",
    "sql_postgres",
    "sql_mysql",
  ]
}

module "postgresql" {
  depends_on = [module.kind]
  source     = "../../modules/postgresql"

  database_name = random_pet.postgres_database.id
  node_ip       = module.kind.node_ip
  node_port     = module.kind.node_ports["postgresql"]
  user_name     = random_pet.postgres_user.id
}

module "mysql" {
  depends_on = [module.kind]
  source     = "../../modules/mysql"

  database_name = random_pet.mysql_database.id
  node_ip       = module.kind.node_ip
  node_port     = module.kind.node_ports["mysql"]
  user_name     = random_pet.mysql_user.id
}

module "redis" {
  depends_on = [module.kind]
  source     = "../../modules/redis"

  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["redis"]
}

module "sql_postgres" {
  depends_on = [module.postgresql, module.redis]
  source     = "../../modules/sql"

  feature   = "postgres"
  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["sql_postgres"]
  secrets = {
    "database_url"      = module.postgresql.cluster_url
    "database_user"     = module.postgresql.user_name
    "database_password" = module.postgresql.user_password
    "database_name"     = module.postgresql.database_name
    "redis_url"         = module.redis.cluster_url
    "redis_password"    = module.redis.redis_password
    "sql_token"         = random_password.sql_postgres_token.result
  }
}

module "sql_mysql" {
  depends_on = [module.mysql, module.redis]
  source     = "../../modules/sql"

  feature   = "mysql"
  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["sql_mysql"]
  secrets = {
    "database_url"      = module.mysql.cluster_url
    "database_user"     = module.mysql.user_name
    "database_password" = module.mysql.user_password
    "database_name"     = module.mysql.database_name
    "redis_password"    = module.redis.redis_password
    "redis_url"         = module.redis.cluster_url
    "redis_password"    = module.redis.redis_password
    "sql_token"         = random_password.sql_mysql_token.result
  }
}
