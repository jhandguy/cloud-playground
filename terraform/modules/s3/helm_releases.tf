resource "helm_release" "s3" {
  name             = "s3"
  namespace        = "s3"
  chart            = "../../../s3/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: ${var.replicas}
    horizontalPodAutoscaler:
      minReplicas: ${var.min_replicas}
      maxReplicas: ${var.max_replicas}
      targetCPUUtilizationPercentage: 100
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
