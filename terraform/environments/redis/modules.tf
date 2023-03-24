module "kind" {
  source = "../../modules/kind"

  cluster_name = "redis"
  node_ports = [
    "postgresql",
    "mysql",
    "redis",
    "sql_postgres_http",
    "sql_postgres_metrics",
    "sql_mysql_http",
    "sql_mysql_metrics",
    "prometheus",
    "alertmanager",
    "grafana",
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
  depends_on = [module.metrics, module.prometheus, module.postgresql, module.redis]
  source     = "../../modules/sql"

  feature            = "postgres"
  node_ip            = module.kind.node_ip
  node_ports         = [module.kind.node_ports["sql_postgres_http"], module.kind.node_ports["sql_postgres_metrics"]]
  prometheus_enabled = true
  secrets = {
    "database_url"      = module.postgresql.cluster_url
    "database_user"     = module.postgresql.user_name
    "database_password" = module.postgresql.user_password
    "database_name"     = module.postgresql.database_name
    "redis_url"         = module.redis.cluster_url
    "redis_password"    = module.redis.redis_password
    "sql_token"         = random_password.sql_postgres_token.result
    "tempo_url"         = module.tempo.otlp_grpc_url
  }
}

module "sql_mysql" {
  depends_on = [module.metrics, module.prometheus, module.mysql, module.redis]
  source     = "../../modules/sql"

  feature            = "mysql"
  node_ip            = module.kind.node_ip
  node_ports         = [module.kind.node_ports["sql_mysql_http"], module.kind.node_ports["sql_mysql_metrics"]]
  prometheus_enabled = true
  secrets = {
    "database_url"      = module.mysql.cluster_url
    "database_user"     = module.mysql.user_name
    "database_password" = module.mysql.user_password
    "database_name"     = module.mysql.database_name
    "redis_password"    = module.redis.redis_password
    "redis_url"         = module.redis.cluster_url
    "redis_password"    = module.redis.redis_password
    "sql_token"         = random_password.sql_mysql_token.result
    "tempo_url"         = module.tempo.otlp_grpc_url
  }
}

module "prometheus" {
  depends_on = [module.kind]
  source     = "../../modules/prometheus"

  alertmanager_node_port = module.kind.node_ports["alertmanager"]
  grafana_dashboards     = ["mysql", "postgres"]
  grafana_datasources    = ["loki", "tempo"]
  grafana_node_port      = module.kind.node_ports["grafana"]
  node_ip                = module.kind.node_ip
  prometheus_node_port   = module.kind.node_ports["prometheus"]
}

module "loki" {
  depends_on = [module.prometheus]
  source     = "../../modules/loki"

  alerting_rules   = ["mysql", "postgres"]
  alertmanager_url = module.prometheus.alertmanager_cluster_url
  labels           = ["level", "message", "target"]
}

module "tempo" {
  depends_on = [module.prometheus]
  source     = "../../modules/tempo"
}

module "metrics" {
  depends_on = [module.kind]
  source     = "../../modules/metrics"
}
