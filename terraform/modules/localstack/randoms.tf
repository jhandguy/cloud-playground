resource "random_id" "buckets" {
  count       = length(var.aws_s3_buckets)
  byte_length = 4
}

resource "random_id" "tables" {
  count       = length(var.aws_dynamo_tables)
  byte_length = 4
}