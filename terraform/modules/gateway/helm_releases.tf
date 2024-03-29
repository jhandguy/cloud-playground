resource "helm_release" "gateway" {
  name             = "gateway"
  namespace        = "gateway"
  chart            = "../../../gateway/helm"
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
%{for name, node_port in var.node_ports~}
      ${name}:
        http: ${node_port.0}
        metrics: ${node_port.1}
%{endfor~}
    ingress:
      host: ${var.ingress_host}
    prometheus:
      enabled: ${var.prometheus_enabled}
      url: http://${var.prometheus_url}
    consul:
      enabled: ${var.consul_enabled}
      ingressGateway:
        port: ${var.ingress_gateway_port}
    csi:
      enabled: ${var.csi_enabled}
      vaultAddress: ${var.vault_url}
    argoRollouts:
      enabled: ${var.argorollouts_enabled}
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
