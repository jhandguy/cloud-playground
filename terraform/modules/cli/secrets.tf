resource "kubernetes_secret" "cli" {
  metadata {
    name      = "cli"
    namespace = kubernetes_namespace.cli.metadata.0.name
  }

  data = {
    gateway_api_key = var.gateway_api_key
  }
}

resource "kubernetes_secret" "cli_image" {
  metadata {
    name      = "cli-image"
    namespace = kubernetes_namespace.cli.metadata.0.name
  }

  data = {
    ".dockerconfigjson" = <<DOCKER
{
  "auths": {
    "${var.image_registry}": {
      "auth": "${base64encode("${var.registry_username}:${var.registry_password}")}"
    }
  }
}
DOCKER
  }

  type = "kubernetes.io/dockerconfigjson"
}