resource "kubernetes_namespace" "cli" {
  metadata {
    name = "cli"
  }
}