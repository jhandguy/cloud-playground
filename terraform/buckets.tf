resource "aws_s3_bucket" "bucket" {
  depends_on = [kubernetes_ingress.localstack]

  bucket = random_pet.bucket.id
}