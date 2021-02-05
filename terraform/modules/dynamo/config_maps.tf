resource "kubernetes_config_map" "dynamo" {
  metadata {
    name      = "dynamo"
    namespace = kubernetes_namespace.dynamo.metadata.0.name
  }

  data = {
    aws_region          = var.aws_region
    aws_dynamo_endpoint = var.aws_dynamo_endpoint
    aws_dynamo_table    = var.aws_dynamo_table
  }
}