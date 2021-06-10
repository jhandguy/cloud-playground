terraform {
  required_version = "~> 1"
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    random = {
      source = "hashicorp/random"
    }
    grafana = {
      source = "grafana/grafana"
    }
  }
}
