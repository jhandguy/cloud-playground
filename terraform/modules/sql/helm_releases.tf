resource "helm_release" "sql" {
  name             = "sql-${var.feature}"
  namespace        = "sql-${var.feature}"
  chart            = "../../../sql/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: ${var.replicas}
    nodePort: ${var.node_port}
    EOF
  ]

  dynamic "set_sensitive" {
    for_each = var.secrets

    content {
      name  = "secrets.${set_sensitive.key}"
      value = base64encode(set_sensitive.value)
    }
  }
}
