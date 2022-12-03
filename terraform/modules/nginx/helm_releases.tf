resource "helm_release" "nginx" {
  name             = "nginx"
  namespace        = "nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  create_namespace = true
  wait             = true
  version          = "4.4.0"

  values = [
    <<-EOF
    controller:
      ingressClassResource:
        default: true
      config:
        ssl-redirect: false
%{if var.prometheus_enabled}
      metrics:
        enabled: true
        serviceMonitor:
          enabled: true
          additionalLabels:
            release: prometheus
%{endif}
      admissionWebhooks:
        enabled: false
      service:
        type: NodePort
        nodePorts:
          http: ${var.node_port}
    defaultBackend:
      enabled: true
    EOF
  ]
}
