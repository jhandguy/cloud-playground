resource "helm_release" "s3" {
  name             = "s3"
  namespace        = "s3"
  chart            = "../s3/helm"
  create_namespace = true
  wait             = true

  values = [<<-EOF
    replicas: 1
    nodePort: ${var.node_port}
    EOF
  ]
}