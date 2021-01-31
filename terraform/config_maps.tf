resource "kubernetes_config_map" "s3" {
  metadata {
    name      = "s3"
    namespace = kubernetes_namespace.s3.metadata.0.name
  }

  data = {
    aws_region      = var.aws_region
    aws_s3_endpoint = "http://${helm_release.localstack.name}.${helm_release.localstack.namespace}.svc.cluster.local:4566"
    aws_s3_bucket   = aws_s3_bucket.s3.id
  }
}

resource "kubernetes_config_map" "dynamo" {
  metadata {
    name      = "dynamo"
    namespace = kubernetes_namespace.dynamo.metadata.0.name
  }

  data = {
    aws_region          = var.aws_region
    aws_dynamo_endpoint = "http://${helm_release.localstack.name}.${helm_release.localstack.namespace}.svc.cluster.local:4566"
    aws_dynamo_table    = aws_dynamodb_table.dynamo.id
  }
}