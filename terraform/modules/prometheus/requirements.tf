terraform {
  required_version = "~> 0.14"
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
