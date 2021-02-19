resource "helm_release" "gateway" {
  name      = "gateway"
  namespace = kubernetes_namespace.gateway.metadata.0.name
  chart     = "../gateway/helm"

  set {
    name  = "replicas"
    value = 1
  }

  set {
    name  = "nodePort"
    value = var.node_port
  }

  set {
    name  = "configMap"
    value = kubernetes_config_map.gateway.metadata.0.name
  }

  set {
    name  = "secret"
    value = kubernetes_secret.gateway.metadata.0.name
  }

  set {
    name  = "image.secret"
    value = kubernetes_secret.gateway_image.metadata.0.name
  }

  set {
    name  = "image.registry"
    value = var.image_registry
  }

  set {
    name  = "image.repository"
    value = var.gateway_image_repository
  }

  set {
    name  = "image.tag"
    value = var.gateway_image_tag
  }
}