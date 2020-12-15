provider "kubernetes" {
  config_context_cluster = "minikube"
}

provider "helm" {
  kubernetes {
    config_context_cluster = "minikube"
  }
}

provider "aws" {
  region                      = "us-east-1"
  access_key                  = "mock_access_key"
  secret_key                  = "mock_secret_key"
  s3_force_path_style         = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3 = "http://${kubernetes_ingress.localstack.load_balancer_ingress.0.ip}"
  }
}