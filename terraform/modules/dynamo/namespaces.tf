resource "kubernetes_namespace" "dynamo" {
  metadata {
    name = "dynamo"
  }
}