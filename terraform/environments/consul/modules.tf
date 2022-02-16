module "minikube" {
  source = "../../modules/minikube"

  node_ports = [
    "localstack",
    "dynamo",
    "s3",
    "gateway_canary",
    "gateway_stable",
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
  aws_s3_buckets    = ["s3"]
  consul_enabled    = true
  node_ip           = var.node_ip
  node_port         = module.minikube.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.metrics, module.vault, module.localstack]
  source     = "../../modules/dynamo"

  consul_enabled     = true
  csi_enabled        = true
  node_ip            = var.node_ip
  node_port          = module.minikube.node_ports["dynamo"]
  prometheus_enabled = true
  vault_url          = module.vault.cluster_url
}

module "s3" {
  depends_on = [module.metrics, module.vault, module.localstack]
  source     = "../../modules/s3"

  consul_enabled     = true
  csi_enabled        = true
  node_ip            = var.node_ip
  node_port          = module.minikube.node_ports["s3"]
  prometheus_enabled = true
  vault_url          = module.vault.cluster_url
}

module "gateway" {
  depends_on = [module.metrics, module.vault, module.dynamo, module.s3]
  source     = "../../modules/gateway"

  consul_enabled       = true
  csi_enabled          = true
  ingress_gateway_port = module.consul.ingress_gateway_port
  ingress_host         = random_pet.gateway_host.id
  node_ip              = var.node_ip
  node_ports = {
    "canary" : module.minikube.node_ports["gateway_canary"],
    "stable" : module.minikube.node_ports["gateway_stable"]
  }
  prometheus_enabled = true
  vault_url          = module.vault.cluster_url
}

module "cli" {
  depends_on = [module.dynamo, module.s3, module.gateway]
  source     = "../../modules/cli"

  csi_enabled = true
  vault_url   = module.vault.cluster_url
}

module "prometheus" {
  source = "../../modules/prometheus"

  alertmanager_node_port = module.minikube.node_ports["alertmanager"]
  grafana_dashboards     = ["dynamo", "s3", "gateway", "cli"]
  grafana_datasources    = ["loki", "tempo"]
  grafana_node_port      = module.minikube.node_ports["grafana"]
  node_ip                = var.node_ip
  prometheus_node_port   = module.minikube.node_ports["prometheus"]
}

module "pushgateway" {
  depends_on = [module.consul]
  source     = "../../modules/pushgateway"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["pushgateway"]
}

module "loki" {
  depends_on = [module.consul]
  source     = "../../modules/loki"

  alerting_rules   = ["dynamo", "s3", "gateway", "cli"]
  alertmanager_url = module.prometheus.alertmanager_url
}

module "tempo" {
  depends_on = [module.consul]
  source     = "../../modules/tempo"

  consul_enabled = true
}

module "metrics" {
  source = "../../modules/metrics"
}

module "csi" {
  source = "../../modules/csi"
}

module "consul" {
  depends_on = [module.prometheus]
  source     = "../../modules/consul"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["consul"]
  node_ports = {
    "gateway" : module.minikube.node_ports["gateway"]
  }
}

module "vault" {
  depends_on = [module.consul, module.csi]
  source     = "../../modules/vault"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["vault"]
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
      "gateway_url"     = module.consul.ingress_gateway_urls["gateway"]
      "gateway_host"    = random_pet.gateway_host.id
      "pushgateway_url" = module.pushgateway.url
      "gateway_token"   = random_password.gateway_token.result
    }
  }
}
