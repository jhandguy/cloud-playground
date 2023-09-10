resource "helm_release" "loki" {
  name             = "loki"
  namespace        = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki-stack"
  create_namespace = true
  wait             = true
  version          = "2.9.11"

  values = [
    <<-EOF
    loki:
      config:
        ruler:
          storage:
            type: local
            local:
              directory: /rules
          rule_path: /tmp/scratch
          alertmanager_url: http://${var.alertmanager_url}
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
    promtail:
      config:
        snippets:
          pipelineStages:
            - cri:
            - json:
                expressions:
%{for label in var.labels~}
                  ${label}:
%{endfor~}
            - labels:
%{for label in var.labels~}
                ${label}:
%{endfor~}
    EOF
  ]
}
