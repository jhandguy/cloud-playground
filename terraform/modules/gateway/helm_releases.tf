resource "helm_release" "gateway" {
  name             = "gateway"
  namespace        = "gateway"
  chart            = "../../../gateway/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: 1
    horizontalPodAutoscaler:
      minReplicas: 1
      maxReplicas: 1
      targetCPUUtilizationPercentage: 50
    services:
%{for name, node_port in var.node_ports~}
      ${name}:
        nodePort: ${node_port}
%{endfor~}
    ingress:
      host: ${var.ingress_host}
    prometheus:
      enabled: ${var.prometheus_enabled}
    consul:
      enabled: ${var.consul_enabled}
      ingressGateway:
        port: ${var.ingress_gateway_port}
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
