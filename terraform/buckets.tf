resource "aws_s3_bucket" "s3" {
  depends_on = [helm_release.localstack]

  bucket = random_pet.s3_bucket.id
}