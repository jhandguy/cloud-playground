resource "kubernetes_namespace" "s3" {
  metadata {
    name = "s3"
  }
}