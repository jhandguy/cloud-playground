resource "aws_s3_bucket" "s3" {
  depends_on = [helm_release.localstack]
  for_each = {
    for index in range(0, length(var.aws_s3_buckets)) : var.aws_s3_buckets[index] => random_id.buckets[index].hex
  }

  bucket        = each.value
  force_destroy = true
}