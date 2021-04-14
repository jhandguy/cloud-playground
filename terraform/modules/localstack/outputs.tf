output "aws_s3_endpoint" {
  value = "${var.node_ip}:${var.node_port}"
}

output "aws_s3_buckets" {
  value = {
    for index in range(0, length(var.aws_s3_buckets)) : var.aws_s3_buckets[index] => aws_s3_bucket.s3[var.aws_s3_buckets[index]].id
  }
  description = "AWS S3 buckets"
}

output "aws_dynamo_endpoint" {
  value = "${var.node_ip}:${var.node_port}"
}

output "aws_dynamo_tables" {
  value = {
    for index in range(0, length(var.aws_s3_buckets)) : var.aws_dynamo_tables[index] => aws_dynamodb_table.dynamo[var.aws_dynamo_tables[index]].id
  }
  description = "AWS DynamoDB tables"
}