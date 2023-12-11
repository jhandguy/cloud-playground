resource "helm_release" "tempo" {
  name             = "tempo"
  namespace        = "tempo"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "tempo"
  create_namespace = true
  wait             = true
  version          = "1.7.1"

  values = [
    <<-EOF
    serviceMonitor:
      enabled: ${var.prometheus_enabled}
      additionalLabels:
        release: prometheus
%{if var.consul_enabled}
    podAnnotations:
      'consul.hashicorp.com/connect-inject': "true"
      'consul.hashicorp.com/connect-service': "tempo"
      'consul.hashicorp.com/connect-service-port': "otlp-grpc"
%{endif}
    EOF
  ]
}
