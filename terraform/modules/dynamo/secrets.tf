resource "kubernetes_secret" "dynamo" {
  metadata {
    name      = "dynamo"
    namespace = kubernetes_namespace.dynamo.metadata.0.name
  }

  data = {
    aws_access_key_id     = var.aws_access_key_id
    aws_secret_access_key = var.aws_secret_access_key
    dynamo_token          = random_password.dynamo_token.result
  }
}

resource "kubernetes_secret" "dynamo_image" {
  metadata {
    name      = "dynamo-image"
    namespace = kubernetes_namespace.dynamo.metadata.0.name
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