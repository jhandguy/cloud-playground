resource "helm_release" "prometheus" {
  name             = "prometheus"
  namespace        = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  create_namespace = true
  wait             = true
  version          = "55.1.0"

  values = [
    <<-EOF
    alertmanager:
      service:
        type: NodePort
        nodePort: ${var.alertmanager_node_port}
      alertmanagerSpec:
        resources:
          requests:
            cpu: 0
            memory: 0
          limits:
            cpu: 0
            memory: 0
    grafana:
      service:
        type: NodePort
        nodePort: ${var.grafana_node_port}
      adminPassword: "${random_password.admin_password.result}"
      additionalDataSources:
%{for datasource in var.grafana_datasources~}
        ${indent(8, file("${path.module}/datasources/${datasource}.yaml"))}
%{endfor~}
      dashboards:
        default:
%{for dashboard in var.grafana_dashboards~}
          ${dashboard}:
            json: |
              ${indent(14, file("${path.module}/dashboards/${dashboard}.json"))}
%{endfor~}
      dashboardProviders:
        dashboardproviders.yaml:
          apiVersion: 1
          providers:
            - name: 'default'
              orgId: 1
              folder: 'Cloud Playground'
              type: file
              disableDeletion: false
              editable: true
              options:
                path: /var/lib/grafana/dashboards/default
    prometheus:
      service:
        type: NodePort
        nodePort: ${var.prometheus_node_port}
      prometheusSpec:
        resources: null
        scrapeInterval: 30s
        evaluationInterval: 30s
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
%{if var.mimir_url != ""}
        remoteWrite:
            - url: http://${var.mimir_url}/api/v1/push
%{endif}
    prometheusOperator:
      tls:
        enabled: false
      admissionWebhooks:
        enabled: false
      prometheusConfigReloader:
        resources:
          requests:
            cpu: 0
            memory: 0
          limits:
            cpu: 0
            memory: 0
    EOF
  ]
}
