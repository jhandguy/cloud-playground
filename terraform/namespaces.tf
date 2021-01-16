resource "kubernetes_namespace" "s3" {
  metadata {
    name = "s3"
  }
}

resource "kubernetes_namespace" "dynamo" {
  metadata {
    name = "dynamo"
  }
}