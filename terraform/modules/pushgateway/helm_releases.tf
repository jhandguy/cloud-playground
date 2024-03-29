resource "helm_release" "pushgateway" {
  name             = "pushgateway"
  namespace        = "pushgateway"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus-pushgateway"
  create_namespace = true
  wait             = true
  version          = "2.4.2"

  values = [
    <<-EOF
    service:
      type: NodePort
      nodePort: ${var.node_port}
    serviceMonitor:
      enabled: ${var.prometheus_enabled}
      namespace: pushgateway
      additionalLabels:
        release: prometheus
%{if var.consul_enabled}
    podAnnotations:
      'consul.hashicorp.com/connect-inject': "true"
      'consul.hashicorp.com/connect-service': "pushgateway"
      'consul.hashicorp.com/connect-service-port': "metrics"
%{endif}
    EOF
  ]
}
