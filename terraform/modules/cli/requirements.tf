terraform {
  required_version = "~> 1"
  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}
