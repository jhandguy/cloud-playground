provider "kubernetes" {
  config_context_cluster = "minikube"
}

provider "helm" {
  kubernetes {
    config_context_cluster = "minikube"
  }
}