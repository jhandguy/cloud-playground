resource "kubernetes_secret" "gateway" {
  metadata {
    name      = "gateway"
    namespace = kubernetes_namespace.gateway.metadata.0.name
  }

  data = {
    gateway_api_key = random_password.gateway_api_key.result
    dynamo_token    = var.dynamo_token
    s3_token        = var.s3_token
  }
}

resource "kubernetes_secret" "gateway_image" {
  metadata {
    name      = "gateway-image"
    namespace = kubernetes_namespace.gateway.metadata.0.name
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