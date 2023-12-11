resource "helm_release" "nginx" {
  name             = "nginx"
  namespace        = "nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  create_namespace = true
  wait             = true
  version          = "4.8.4"

  values = [
    <<-EOF
    controller:
      ingressClassResource:
        default: true
      config:
        ssl-redirect: false
      metrics:
        enabled: true
        serviceMonitor:
          enabled: ${var.prometheus_enabled}
          additionalLabels:
            release: prometheus
      admissionWebhooks:
        enabled: false
      service:
        type: NodePort
        nodePorts:
          http: ${var.node_ports.0}
          https: ${var.node_ports.1}
      resources: null
    defaultBackend:
      enabled: true
    EOF
  ]
}
