resource "helm_release" "loki" {
  name             = "loki"
  namespace        = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki-stack"
  create_namespace = true
  wait             = true
  version          = "2.4.1"

  values = [<<-EOF
    loki:
      service:
        type: NodePort
        nodePort: ${var.node_port}
      serviceMonitor:
        enabled: true
        additionalLabels:
          release: prometheus
    EOF
  ]
}