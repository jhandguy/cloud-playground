resource "helm_release" "cli" {
  name      = "cli"
  namespace = kubernetes_namespace.cli.metadata.0.name
  chart     = "../cli/helm"
  wait      = true

  set {
    name  = "configMap"
    value = kubernetes_config_map.cli.metadata.0.name
  }

  set {
    name  = "secret"
    value = kubernetes_secret.cli.metadata.0.name
  }

  set {
    name  = "image.secret"
    value = kubernetes_secret.cli_image.metadata.0.name
  }

  set {
    name  = "image.registry"
    value = var.image_registry
  }

  set {
    name  = "image.repository"
    value = var.cli_image_repository
  }

  set {
    name  = "image.tag"
    value = var.cli_image_tag
  }

  set {
    name  = "test.rounds"
    value = 100
  }
}