terraform {
  required_version = "~> 0.15"
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
