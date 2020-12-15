resource "aws_s3_bucket" "bucket" {
  depends_on = [kubernetes_ingress.localstack]

  acl    = "public-read"
}