module "kind" {
  source = "../../modules/kind"

  cluster_name = "consul"
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
    "gateway",
    "prometheus",
    "alertmanager",
    "grafana",
    "pushgateway",
    "consul",
    "vault",
  ]
}

module "localstack" {
  depends_on = [module.consul]
  source     = "../../modules/localstack"

  aws_dynamo_tables = ["dynamo"]
  aws_s3_buckets    = ["mimir", "s3"]
  consul_enabled    = true
  node_ip           = module.kind.node_ip
  node_port         = module.kind.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.metrics, module.vault, module.localstack]
  source     = "../../modules/dynamo"

  max_replicas   = 2
  consul_enabled = true
  csi_enabled    = true
  node_ip        = module.kind.node_ip
  node_ports     = [module.kind.node_ports["dynamo_grpc"], module.kind.node_ports["dynamo_metrics"]]
  vault_url      = module.vault.cluster_url
}

module "s3" {
  depends_on = [module.metrics, module.vault, module.localstack]
  source     = "../../modules/s3"

  max_replicas   = 2
  consul_enabled = true
  csi_enabled    = true
  node_ip        = module.kind.node_ip
  node_ports     = [module.kind.node_ports["s3_grpc"], module.kind.node_ports["s3_metrics"]]
  vault_url      = module.vault.cluster_url
}

module "gateway" {
  depends_on = [module.dynamo, module.s3]
  source     = "../../modules/gateway"

  max_replicas         = 2
  consul_enabled       = true
  csi_enabled          = true
  ingress_gateway_port = module.consul.ingress_gateway_port
  ingress_host         = random_pet.gateway_host.id
  node_ip              = module.kind.node_ip
  node_ports = {
    "canary" : [module.kind.node_ports["gateway_canary_http"], module.kind.node_ports["gateway_canary_metrics"]],
    "stable" : [module.kind.node_ports["gateway_stable_http"], module.kind.node_ports["gateway_stable_metrics"]]
  }
  vault_url = module.vault.cluster_url
}

module "cli" {
  depends_on = [module.dynamo, module.s3, module.gateway]
  source     = "../../modules/cli"

  csi_enabled = true
  vault_url   = module.vault.cluster_url
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
  depends_on = [module.consul, module.prometheus]
  source     = "../../modules/pushgateway"

  consul_enabled = true
  node_ip        = module.kind.node_ip
  node_port      = module.kind.node_ports["pushgateway"]
}

module "loki" {
  depends_on = [module.consul, module.prometheus]
  source     = "../../modules/loki"

  alerting_rules   = ["dynamo", "s3", "gateway", "cli"]
  alertmanager_url = module.prometheus.alertmanager_cluster_url
}

module "tempo" {
  depends_on = [module.consul, module.prometheus]
  source     = "../../modules/tempo"

  consul_enabled = true
}

module "metrics" {
  depends_on = [module.kind]
  source     = "../../modules/metrics"
}

module "csi" {
  depends_on = [module.kind]
  source     = "../../modules/csi"
}

module "consul" {
  depends_on = [module.kind]
  source     = "../../modules/consul"

  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["consul"]
  node_ports = {
    "gateway" : module.kind.node_ports["gateway"]
  }
}

module "vault" {
  depends_on = [module.csi]
  source     = "../../modules/vault"

  node_ip   = module.kind.node_ip
  node_port = module.kind.node_ports["vault"]
  secrets = {
    "s3" : {
      "aws_region"            = var.aws_region
      "aws_access_key_id"     = var.aws_access_key_id
      "aws_secret_access_key" = var.aws_secret_access_key
      "aws_s3_bucket"         = module.localstack.aws_s3_buckets["s3"]
      "s3_token"              = random_password.s3_token.result
    },
    "dynamo" : {
      "aws_region"            = var.aws_region
      "aws_access_key_id"     = var.aws_access_key_id
      "aws_secret_access_key" = var.aws_secret_access_key
      "aws_dynamo_table"      = module.localstack.aws_dynamo_tables["dynamo"]
      "dynamo_token"          = random_password.dynamo_token.result
    },
    "gateway" : {
      "gateway_token" = random_password.gateway_token.result
      "dynamo_token"  = random_password.dynamo_token.result
      "s3_token"      = random_password.s3_token.result
    },
    "cli" : {
      "gateway_url"     = module.consul.ingress_gateway_cluster_url
      "gateway_host"    = random_pet.gateway_host.id
      "pushgateway_url" = module.pushgateway.cluster_url
      "gateway_token"   = random_password.gateway_token.result
    }
  }
}
