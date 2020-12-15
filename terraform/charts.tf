resource "helm_release" "jenkins" {
  name             = "jenkins"
  namespace        = "jenkins"
  repository       = "https://charts.jenkins.io"
  chart            = "jenkins"
  create_namespace = true

  set {
    name  = "controller.jenkinsUriPrefix"
    value = "/jenkins"
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
    value = "s3"
  }
}