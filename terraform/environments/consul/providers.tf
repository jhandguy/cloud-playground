provider "helm" {
  kubernetes {
    config_context_cluster = module.kind.cluster_context
    config_path            = "~/.kube/config"
  }
}

provider "aws" {
  region                      = var.aws_region
  access_key                  = var.aws_access_key_id
  secret_key                  = var.aws_secret_access_key
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3       = "http://${module.kind.node_ip}:${module.kind.node_ports["localstack"]}"
    dynamodb = "http://${module.kind.node_ip}:${module.kind.node_ports["localstack"]}"
  }
}
