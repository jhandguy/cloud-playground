resource "helm_release" "haproxy" {
  name             = "haproxy"
  namespace        = "haproxy"
  repository       = "https://haproxytech.github.io/helm-charts"
  chart            = "kubernetes-ingress"
  create_namespace = true
  wait             = true
  version          = "1.33.0"

  values = [
    <<-EOF
    controller:
      replicaCount: 1
      ingressClassResource:
        default: true
      config:
        ssl-redirect: "false"
      serviceMonitor:
        enabled: ${var.prometheus_enabled}
        extraLabels:
          release: prometheus
      service:
        type: NodePort
        nodePorts:
          http: ${var.node_ports.0}
          https: ${var.node_ports.1}
          stat: ${var.node_ports.2}
          prometheus: ${var.node_ports.3}
    defaultBackend:
      enabled: true
    EOF
  ]
}
