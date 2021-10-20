resource "helm_release" "loki" {
  name             = "loki"
  namespace        = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki-stack"
  create_namespace = true
  wait             = true
  version          = "2.4.1"

  values = [
    <<-EOF
    loki:
      service:
        type: NodePort
        nodePort: ${var.node_port}
      config:
        ruler:
          storage:
            type: local
            local:
              directory: /rules
          rule_path: /tmp/scratch
          alertmanager_url: http://${var.node_ip}:${var.alertmanager_node_port}
          ring:
            kvstore:
              store: inmemory
          enable_api: true
      serviceMonitor:
        enabled: true
        additionalLabels:
          release: prometheus
      alerting_groups:
%{for rule in var.alerting_rules~}
        ${indent(8, file("${path.module}/rules/${rule}.yaml"))}
%{endfor~}
    EOF
  ]
}
