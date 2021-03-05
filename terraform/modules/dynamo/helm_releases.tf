resource "helm_release" "dynamo" {
  name      = "dynamo"
  namespace = kubernetes_namespace.dynamo.metadata.0.name
  chart     = "../dynamo/helm"
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
    value = kubernetes_config_map.dynamo.metadata.0.name
  }

  set {
    name  = "secret"
    value = kubernetes_secret.dynamo.metadata.0.name
  }

  set {
    name  = "image.secret"
    value = kubernetes_secret.dynamo_image.metadata.0.name
  }

  set {
    name  = "image.registry"
    value = var.image_registry
  }

  set {
    name  = "image.repository"
    value = var.dynamo_image_repository
  }

  set {
    name  = "image.tag"
    value = var.dynamo_image_tag
  }
}