terraform {
  required_version = "~> 0.15"
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
    helm = {
      source = "hashicorp/helm"
    }
    aws = {
      source = "hashicorp/aws"
    }
    random = {
      source = "hashicorp/random"
    }
    grafana = {
      source = "grafana/grafana"
    }
  }
}
