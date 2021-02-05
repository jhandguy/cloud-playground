resource "helm_release" "localstack" {
  name             = "localstack"
  namespace        = "localstack"
  repository       = "http://helm.localstack.cloud"
  chart            = "localstack"
  create_namespace = true

  set {
    name  = "startServices"
    value = "s3\\,dynamodb"
  }

  set {
    name  = "nodePorts.edgePort"
    value = var.node_port
  }
}