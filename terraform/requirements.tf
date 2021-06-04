terraform {
  required_version = "~> 0.15"
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    aws = {
      source = "hashicorp/aws"
    }
    random = {
      source = "hashicorp/random"
    }
    null = {
      source = "hashicorp/null"
    }
    grafana = {
      source = "grafana/grafana"
    }
  }
}
