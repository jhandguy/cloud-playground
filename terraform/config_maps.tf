resource "kubernetes_config_map" "s3" {
  metadata {
    name      = "s3"
    namespace = kubernetes_namespace.s3.metadata.0.name
  }

  data = {
    aws_region      = var.aws_region
    aws_s3_endpoint = "http://${kubernetes_ingress.localstack.load_balancer_ingress.0.ip}/"
    aws_s3_bucket   = aws_s3_bucket.bucket.id
  }
}