resource "helm_release" "prometheus" {
  name             = "prometheus"
  namespace        = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  create_namespace = true

  values = [<<-EOF
    alertmanager:
      service:
        type: NodePort
        nodePort: ${var.alertmanager_node_port}
    grafana:
      service:
        type: NodePort
        nodePort: ${var.grafana_node_port}
      adminPassword: "${random_password.admin_password.result}"
    prometheus:
      service:
        type: NodePort
        nodePort: ${var.prometheus_node_port}
    prometheusOperator:
      tls:
        enabled: false
      admissionWebhooks:
        enabled: false
    EOF
  ]
}