resource "kubernetes_secret" "s3" {
  metadata {
    name      = "s3"
    namespace = kubernetes_namespace.s3.metadata.0.name
  }

  data = {
    aws_access_key_id     = var.aws_access_key_id
    aws_secret_access_key = var.aws_secret_access_key
    api_key               = random_password.s3_api_key.result
  }
}

resource "kubernetes_secret" "s3_image" {
  metadata {
    name      = "s3-image"
    namespace = kubernetes_namespace.s3.metadata.0.name
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

resource "kubernetes_secret" "dynamo" {
  metadata {
    name      = "dynamo"
    namespace = kubernetes_namespace.dynamo.metadata.0.name
  }

  data = {
    aws_access_key_id     = var.aws_access_key_id
    aws_secret_access_key = var.aws_secret_access_key
    api_key               = random_password.dynamo_api_key.result
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