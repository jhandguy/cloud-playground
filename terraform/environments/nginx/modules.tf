module "kind" {
  source = "../../modules/kind"

  cluster_name = var.argorollouts_enabled ? "argorollouts" : "nginx"
  node_ports = [
    "localstack",
    "dynamo_grpc",
    "dynamo_metrics",
    "s3_grpc",
    "s3_metrics",
    "gateway_canary_http",
    "gateway_canary_metrics",
    "gateway_stable_http",
    "gateway_stable_metrics",
    "prometheus",
    "alertmanager",
    "grafana",
    "pushgateway",
    "nginx_http",
    "nginx_https",
    "argorollouts",
  ]
}

module "localstack" {
  depends_on = [module.kind]
  source     = "../../modules/localstack"

  aws_dynamo_tables = ["dynamo"]
  aws_s3_buckets    = ["mimir", "s3"]
  node_ip           = module.kind.node_ip
  node_port         = module.kind.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.metrics, module.prometheus, module.localstack]
  source     = "../../modules/dynamo"

  max_replicas = 2
  node_ip      = module.kind.node_ip
  node_ports   = [module.kind.node_ports["dynamo_grpc"], module.kind.node_ports["dynamo_metrics"]]
  secrets = {
    "aws_region"            = var.aws_region
    "aws_access_key_id"     = var.aws_access_key_id
    "aws_secret_access_key" = var.aws_secret_access_key
    "aws_dynamo_endpoint"   = module.localstack.aws_dynamo_cluster_endpoint
    "aws_dynamo_table"      = module.localstack.aws_dynamo_tables["dynamo"]
    "dynamo_token"          = random_password.dynamo_token.result
    "tempo_url"             = module.tempo.otlp_grpc_url
  }
}

module "s3" {
  depends_on = [module.metrics, module.prometheus, module.localstack]
  source     = "../../modules/s3"

  max_replicas = 2
  node_ip      = module.kind.node_ip
  node_ports   = [module.kind.node_ports["s3_grpc"], module.kind.node_ports["s3_metrics"]]
  secrets = {
    "aws_region"            = var.aws_region
    "aws_access_key_id"     = var.aws_access_key_id
    "aws_secret_access_key" = var.aws_secret_access_key
    "aws_s3_endpoint"       = module.localstack.aws_s3_cluster_endpoint
    "aws_s3_bucket"         = module.localstack.aws_s3_buckets["s3"]
    "s3_token"              = random_password.s3_token.result
    "tempo_url"             = module.tempo.otlp_grpc_url
  }
}

module "gateway" {
  depends_on = [module.argorollouts, module.dynamo, module.s3]
  source     = "../../modules/gateway"

  argorollouts_enabled = var.argorollouts_enabled
  ingress_host         = random_pet.gateway_host.id
  min_replicas         = var.argorollouts_enabled ? 2 : 1
  max_replicas         = var.argorollouts_enabled ? 4 : 2
  node_ip              = module.kind.node_ip
  node_ports = {
    "canary" : [module.kind.node_ports["gateway_canary_http"], module.kind.node_ports["gateway_canary_metrics"]],
    "stable" : [module.kind.node_ports["gateway_stable_http"], module.kind.node_ports["gateway_stable_metrics"]]
  }
  prometheus_url = module.prometheus.prometheus_cluster_url
  replicas       = var.argorollouts_enabled ? 2 : 1
  secrets = {
    "gateway_token" = random_password.gateway_token.result
    "dynamo_url"    = module.dynamo.cluster_url
    "dynamo_token"  = random_password.dynamo_token.result
    "s3_url"        = module.s3.cluster_url
    "s3_token"      = random_password.s3_token.result
    "tempo_url"     = module.tempo.otlp_grpc_url
  }
}

module "cli" {
  depends_on = [module.dynamo, module.s3, module.gateway]
  source     = "../../modules/cli"

  secrets = {
    "gateway_url"     = module.nginx.cluster_url
    "gateway_host"    = module.gateway.host
    "pushgateway_url" = module.pushgateway.cluster_url
    "gateway_token"   = random_password.gateway_token.result
  }
}

module "mimir" {
  depends_on = [module.localstack]
  source     = "../../modules/mimir"

  aws_access_key_id       = var.aws_access_key_id
  aws_region              = var.aws_region
  aws_secret_access_key   = var.aws_secret_access_key
  aws_s3_bucket           = module.localstack.aws_s3_buckets["mimir"]
  aws_s3_cluster_endpoint = module.localstack.aws_s3_cluster_endpoint
}

module "prometheus" {
  depends_on = [module.mimir]
  source     = "../../modules/prometheus"

  alertmanager_node_port = module.kind.node_ports["alertmanager"]
  grafana_dashboards     = ["dynamo", "s3", "gateway", "cli"]
  grafana_datasources    = ["loki", "mimir", "tempo"]
  grafana_node_port      = module.kind.node_ports["grafana"]
  mimir_url              = module.mimir.cluster_url
  node_ip                = module.kind.node_ip
  prometheus_node_port   = module.kind.node_ports["prometheus"]
}

module "pushgateway" {
  depends_on = [module.prometheus]
  source     = "../../modules/pushgateway"

  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["pushgateway"]
}

module "loki" {
  depends_on = [module.prometheus]
  source     = "../../modules/loki"

  alerting_rules   = ["dynamo", "s3", "gateway", "cli"]
  alertmanager_url = module.prometheus.alertmanager_cluster_url
}

module "tempo" {
  depends_on = [module.prometheus]
  source     = "../../modules/tempo"
}

module "metrics" {
  depends_on = [module.kind]
  source     = "../../modules/metrics"
}

module "nginx" {
  depends_on = [module.prometheus]
  source     = "../../modules/nginx"

  node_ip    = module.kind.node_ip
  node_ports = [module.kind.node_ports["nginx_http"], module.kind.node_ports["nginx_https"]]
}

module "certmanager" {
  depends_on = [module.kind]
  source     = "../../modules/certmanager"
}

module "argorollouts" {
  count = var.argorollouts_enabled ? 1 : 0

  depends_on = [module.prometheus, module.nginx]
  source     = "../../modules/argorollouts"

  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["argorollouts"]
}
