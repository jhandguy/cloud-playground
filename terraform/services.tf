data "kubernetes_service" "jenkins" {
  metadata {
    name      = helm_release.jenkins.name
    namespace = helm_release.jenkins.namespace
  }
}

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