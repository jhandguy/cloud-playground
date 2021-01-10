resource "helm_release" "jenkins" {
  name             = "jenkins"
  namespace        = "jenkins"
  repository       = "https://charts.jenkins.io"
  chart            = "jenkins"
  create_namespace = true

  set {
    name  = "controller.jenkinsUriPrefix"
    value = "/${random_pet.jenkins_uri_prefix.id}"
  }
}

resource "helm_release" "localstack" {
  name             = "localstack"
  namespace        = "localstack"
  repository       = "http://helm.localstack.cloud"
  chart            = "localstack"
  create_namespace = true

  set {
    name  = "startServices"
    value = "s3\\,dynamodb"
  }
}

resource "helm_release" "s3" {
  name      = "s3"
  namespace = kubernetes_namespace.s3.metadata.0.name
  chart     = "../s3/helm"

  set {
    name  = "replicas"
    value = 1
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

  set {
    name  = "uriPrefix"
    value = "/${random_pet.s3_uri_prefix.id}"
  }

  set {
    name  = "healthPath"
    value = "/${random_pet.s3_health_path.id}"
  }
}

resource "helm_release" "dynamo" {
  name      = "dynamo"
  namespace = kubernetes_namespace.dynamo.metadata.0.name
  chart     = "../dynamo/helm"

  set {
    name  = "replicas"
    value = 1
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

  set {
    name  = "uriPrefix"
    value = "/${random_pet.dynamo_uri_prefix.id}"
  }

  set {
    name  = "healthPath"
    value = "/${random_pet.dynamo_health_path.id}"
  }
}