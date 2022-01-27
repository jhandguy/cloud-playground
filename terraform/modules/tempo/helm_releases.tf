resource "helm_release" "tempo" {
  name             = "tempo"
  namespace        = "tempo"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "tempo"
  create_namespace = true
  wait             = true
  version          = "0.13.0"

  values = [
    <<-EOF
    serviceMonitor:
      enabled: true
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
