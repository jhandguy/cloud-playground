module "minikube" {
  source = "./modules/minikube"

  node_ports = ["localstack", "dynamo", "s3", "gateway", "prometheus", "alertmanager", "grafana"]
}

module "localstack" {
  source = "./modules/localstack"

  aws_access_key_id     = var.aws_access_key_id
  aws_dynamo_tables     = ["dynamo"]
  aws_region            = var.aws_region
  aws_s3_buckets        = ["s3"]
  aws_secret_access_key = var.aws_secret_access_key
  node_ip               = var.node_ip
  node_port             = module.minikube.node_ports["localstack"]
}

module "dynamo" {
  depends_on = [module.prometheus]
  source     = "./modules/dynamo"

  aws_access_key_id       = var.aws_access_key_id
  aws_dynamo_endpoint     = module.localstack.localstack_endpoint
  aws_dynamo_table        = module.localstack.aws_dynamo_tables["dynamo"]
  aws_region              = var.aws_region
  aws_secret_access_key   = var.aws_secret_access_key
  dynamo_image_repository = var.dynamo_image_repository
  dynamo_image_tag        = var.dynamo_image_tag
  image_registry          = var.image_registry
  node_ip                 = var.node_ip
  node_port               = module.minikube.node_ports["dynamo"]
  registry_password       = var.registry_password
  registry_username       = var.registry_username
}

module "s3" {
  depends_on = [module.prometheus]
  source     = "./modules/s3"

  aws_access_key_id     = var.aws_access_key_id
  aws_region            = var.aws_region
  aws_s3_bucket         = module.localstack.aws_s3_buckets["s3"]
  aws_s3_endpoint       = module.localstack.localstack_endpoint
  aws_secret_access_key = var.aws_secret_access_key
  image_registry        = var.image_registry
  node_ip               = var.node_ip
  node_port             = module.minikube.node_ports["s3"]
  registry_password     = var.registry_password
  registry_username     = var.registry_username
  s3_image_repository   = var.s3_image_repository
  s3_image_tag          = var.s3_image_tag
}

module "gateway" {
  depends_on = [module.prometheus]
  source     = "./modules/gateway"

  dynamo_token             = module.dynamo.token
  dynamo_url               = module.dynamo.url
  gateway_image_repository = var.gateway_image_repository
  gateway_image_tag        = var.gateway_image_tag
  image_registry           = var.image_registry
  node_ip                  = var.node_ip
  node_port                = module.minikube.node_ports["gateway"]
  registry_password        = var.registry_password
  registry_username        = var.registry_username
  s3_token                 = module.s3.token
  s3_url                   = module.s3.url
}

module "prometheus" {
  source = "./modules/prometheus"

  alertmanager_node_port = module.minikube.node_ports["alertmanager"]
  grafana_dashboards     = ["dynamo", "s3"]
  grafana_node_port      = module.minikube.node_ports["grafana"]
  node_ip                = var.node_ip
  prometheus_node_port   = module.minikube.node_ports["prometheus"]
}