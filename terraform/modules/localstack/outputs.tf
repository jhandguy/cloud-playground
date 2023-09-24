output "aws_s3_endpoint" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "AWS S3 endpoint"
}

output "aws_s3_cluster_endpoint" {
  value       = "${helm_release.localstack.name}.${helm_release.localstack.namespace}.svc.cluster.local:4566"
  description = "AWS S3 cluster endpoint"
}

output "aws_s3_buckets" {
  value = {
    for index in range(length(var.aws_s3_buckets)) : var.aws_s3_buckets[index] => aws_s3_bucket.s3[var.aws_s3_buckets[index]].id
  }
  description = "AWS S3 buckets"
}

output "aws_dynamo_endpoint" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "AWS DynamoDB endpoint"
}

output "aws_dynamo_cluster_endpoint" {
  value       = "${helm_release.localstack.name}.${helm_release.localstack.namespace}.svc.cluster.local:4566"
  description = "AWS DynamoDB cluster endpoint"
}

output "aws_dynamo_tables" {
  value = {
    for index in range(length(var.aws_dynamo_tables)) : var.aws_dynamo_tables[index] => aws_dynamodb_table.dynamo[var.aws_dynamo_tables[index]].id
  }
  description = "AWS DynamoDB tables"
}
