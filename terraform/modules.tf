module "minikube" {
  source = "./modules/minikube"

  node_ports = ["localstack", "dynamo", "s3", "canary_gateway", "prod_gateway", "gateway", "prometheus", "alertmanager", "grafana", "pushgateway"]
}

module "localstack" {
  depends_on = [module.consul]
  source     = "./modules/localstack"

  aws_access_key_id     = var.aws_access_key_id
  aws_dynamo_tables     = ["dynamo"]
  aws_region            = var.aws_region
  aws_s3_buckets        = ["s3"]
  aws_secret_access_key = var.aws_secret_access_key
  node_ip               = var.node_ip
  node_port             = module.minikube.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.localstack]
  source     = "./modules/dynamo"

  aws_access_key_id     = var.aws_access_key_id
  aws_dynamo_table      = module.localstack.aws_dynamo_tables["dynamo"]
  aws_region            = var.aws_region
  aws_secret_access_key = var.aws_secret_access_key
  image_registry        = var.image_registry
  node_ip               = var.node_ip
  node_port             = module.minikube.node_ports["dynamo"]
  registry_password     = var.registry_password
  registry_username     = var.registry_username
}

module "s3" {
  depends_on = [module.localstack]
  source     = "./modules/s3"

  aws_access_key_id     = var.aws_access_key_id
  aws_region            = var.aws_region
  aws_s3_bucket         = module.localstack.aws_s3_buckets["s3"]
  aws_secret_access_key = var.aws_secret_access_key
  image_registry        = var.image_registry
  node_ip               = var.node_ip
  node_port             = module.minikube.node_ports["s3"]
  registry_password     = var.registry_password
  registry_username     = var.registry_username
}

module "gateway" {
  depends_on = [module.consul]
  source     = "./modules/gateway"

  dynamo_token         = module.dynamo.token
  image_registry       = var.image_registry
  ingress_gateway_port = module.consul.ingress_gateway_port
  node_ip              = var.node_ip
  node_ports = {
    "canary" : module.minikube.node_ports["canary_gateway"],
    "prod" : module.minikube.node_ports["prod_gateway"]
  }
  registry_password = var.registry_password
  registry_username = var.registry_username
  s3_token          = module.s3.token
}

module "prometheus" {
  source = "./modules/prometheus"

  alertmanager_node_port = module.minikube.node_ports["alertmanager"]
  grafana_dashboards     = ["dynamo", "s3", "gateway", "cli"]
  grafana_node_port      = module.minikube.node_ports["grafana"]
  node_ip                = var.node_ip
  prometheus_node_port   = module.minikube.node_ports["prometheus"]
}

module "pushgateway" {
  depends_on = [module.consul]
  source     = "./modules/pushgateway"

  node_ip   = var.node_ip
  node_port = module.minikube.node_ports["pushgateway"]
}

module "cli" {
  depends_on = [module.dynamo, module.s3, module.gateway]
  source     = "./modules/cli"

  gateway_api_key   = module.gateway.api_key
  gateway_url       = module.consul.ingress_gateway_urls["gateway"]
  image_registry    = var.image_registry
  pushgateway_url   = module.pushgateway.url
  registry_password = var.registry_password
  registry_username = var.registry_username
}

module "consul" {
  depends_on = [module.prometheus]
  source     = "./modules/consul"

  node_ip = var.node_ip
  node_ports = {
    "gateway" : module.minikube.node_ports["gateway"]
  }
}