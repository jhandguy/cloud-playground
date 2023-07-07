provider "helm" {
  kubernetes {
    config_context_cluster = module.kind.cluster_context
    config_path            = "~/.kube/config"
  }
}
