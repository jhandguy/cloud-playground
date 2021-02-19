resource "kubernetes_namespace" "gateway" {
  metadata {
    name = "gateway"
  }
}