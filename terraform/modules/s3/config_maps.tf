resource "kubernetes_config_map" "s3" {
  metadata {
    name      = "s3"
    namespace = kubernetes_namespace.s3.metadata.0.name
  }

  data = {
    aws_region      = var.aws_region
    aws_s3_endpoint = var.aws_s3_endpoint
    aws_s3_bucket   = var.aws_s3_bucket
  }
}