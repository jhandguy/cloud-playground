resource "helm_release" "dynamo" {
  name             = "dynamo"
  namespace        = "dynamo"
  chart            = "../dynamo/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [<<-EOF
    replicas: 1
    nodePort: ${var.node_port}
    EOF
  ]
}