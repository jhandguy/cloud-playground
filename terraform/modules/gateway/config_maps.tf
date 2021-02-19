resource "kubernetes_config_map" "gateway" {
  metadata {
    name      = "gateway"
    namespace = kubernetes_namespace.gateway.metadata.0.name
  }

  data = {
    dynamo_url = var.dynamo_url
    s3_url     = var.s3_url
  }
}