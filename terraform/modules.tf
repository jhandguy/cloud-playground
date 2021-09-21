module "minikube" {
  source = "./modules/minikube"

  node_ports = [
    "localstack",
    "dynamo",
    "s3",
    "gateway_canary",
    "gateway_prod",
    "gateway",
    "prometheus",
    "alertmanager",
    "grafana",
    "loki",
    "pushgateway",
    "consul",
    "vault",
  ]
}

module "localstack" {
  depends_on = [module.consul]
  source     = "./modules/localstack"

  aws_dynamo_tables = ["dynamo"]
  aws_s3_buckets    = ["s3"]
  node_ip           = var.node_ip
  node_port         = module.minikube.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.vault, module.localstack]
  source     = "./modules/dynamo"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["dynamo"]
}

module "s3" {
  depends_on = [module.vault, module.localstack]
  source     = "./modules/s3"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["s3"]
}

module "gateway" {
  depends_on = [module.vault, module.dynamo, module.s3]
  source     = "./modules/gateway"

  ingress_gateway_port = module.consul.ingress_gateway_port
  node_ip              = var.node_ip
  node_ports = {
    "canary" : module.minikube.node_ports["gateway_canary"],
    "prod" : module.minikube.node_ports["gateway_prod"]
  }
}

module "prometheus" {
  depends_on = [module.metrics]
  source     = "./modules/prometheus"

  alertmanager_node_port = module.minikube.node_ports["alertmanager"]
  grafana_dashboards     = ["dynamo", "s3", "gateway", "cli"]
  grafana_node_port      = module.minikube.node_ports["grafana"]
  loki_node_port         = module.minikube.node_ports["loki"]
  node_ip                = var.node_ip
  prometheus_node_port   = module.minikube.node_ports["prometheus"]
}

module "pushgateway" {
  depends_on = [module.consul]
  source     = "./modules/pushgateway"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["pushgateway"]
}

module "loki" {
  depends_on = [module.prometheus]
  source     = "./modules/loki"

  node_port = module.minikube.node_ports["loki"]
}

module "cli" {
  depends_on = [module.dynamo, module.s3, module.gateway]
  source     = "./modules/cli"
}

module "consul" {
  depends_on = [module.prometheus]
  source     = "./modules/consul"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["consul"]
  node_ports = {
    "gateway" : module.minikube.node_ports["gateway"]
  }
}

module "csi" {
  source = "./modules/csi"
}

module "metrics" {
  source = "./modules/metrics"
}

module "vault" {
  depends_on = [module.consul, module.csi]
  source     = "./modules/vault"

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
      "pushgateway_url" = module.pushgateway.url
      "gateway_token"   = random_password.gateway_token.result
    }
  }
}