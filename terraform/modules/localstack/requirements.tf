terraform {
  required_version = "~> 1"
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
  }
}
