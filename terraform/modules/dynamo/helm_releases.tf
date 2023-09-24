resource "helm_release" "dynamo" {
  name             = "dynamo"
  namespace        = "dynamo"
  chart            = "../../../dynamo/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: ${var.replicas}
    horizontalPodAutoscaler:
      minReplicas: ${var.min_replicas}
      maxReplicas: ${var.max_replicas}
    nodePorts:
      grpc: ${var.node_ports.0}
      metrics: ${var.node_ports.1}
    prometheus:
      enabled: ${var.prometheus_enabled}
    consul:
      enabled: ${var.consul_enabled}
    csi:
      enabled: ${var.csi_enabled}
      vaultAddress: ${var.vault_url}
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
