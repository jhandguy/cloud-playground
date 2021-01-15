data "kubernetes_service" "localstack" {
  metadata {
    name      = helm_release.localstack.name
    namespace = helm_release.localstack.namespace
  }
}

data "kubernetes_service" "s3" {
  metadata {
    name      = helm_release.s3.name
    namespace = helm_release.s3.namespace
  }
}

data "kubernetes_service" "dynamo" {
  metadata {
    name      = helm_release.dynamo.name
    namespace = helm_release.dynamo.namespace
  }
}