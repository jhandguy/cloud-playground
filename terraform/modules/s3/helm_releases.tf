resource "helm_release" "s3" {
  name      = "s3"
  namespace = kubernetes_namespace.s3.metadata.0.name
  chart     = "../s3/helm"
  wait      = true

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
    value = kubernetes_config_map.s3.metadata.0.name
  }

  set {
    name  = "secret"
    value = kubernetes_secret.s3.metadata.0.name
  }

  set {
    name  = "image.secret"
    value = kubernetes_secret.s3_image.metadata.0.name
  }

  set {
    name  = "image.registry"
    value = var.image_registry
  }

  set {
    name  = "image.repository"
    value = var.s3_image_repository
  }

  set {
    name  = "image.tag"
    value = var.s3_image_tag
  }
}