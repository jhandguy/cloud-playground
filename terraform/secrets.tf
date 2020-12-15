data "kubernetes_secret" "jenkins" {
  metadata {
    name      = helm_release.jenkins.name
    namespace = helm_release.jenkins.namespace
  }
}