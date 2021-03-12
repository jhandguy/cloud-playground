resource "helm_release" "prometheus" {
  name             = "prometheus"
  namespace        = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  create_namespace = true
  wait             = true

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
    kubeControllerManager:
      enabled: false
    kubeEtcd:
      enabled: false
    kubeScheduler:
      enabled: false
    prometheus:
      service:
        type: NodePort
        nodePort: ${var.prometheus_node_port}
      prometheusSpec:
        scrapeInterval: 1s
        evaluationInterval: 1s
        ruleSelector:
          matchExpressions:
            - key: release
              operator: In
              values:
                - prometheus
                - pushgateway
        serviceMonitorSelector:
          matchExpressions:
            - key: release
              operator: In
              values:
                - prometheus
                - pushgateway
    prometheusOperator:
      tls:
        enabled: false
      admissionWebhooks:
        enabled: false
    EOF
  ]
}

resource "helm_release" "pushgateway" {
  depends_on = [helm_release.prometheus]

  name             = "pushgateway"
  namespace        = "pushgateway"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus-pushgateway"
  create_namespace = true
  wait             = true

  values = [<<-EOF
    service:
      type: NodePort
      nodePort: ${var.pushgateway_node_port}
    serviceMonitor:
      enabled: true
      namespace: pushgateway
    EOF
  ]
}