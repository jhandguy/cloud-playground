resource "helm_release" "argorollouts" {
  name             = "argorollouts"
  namespace        = "argorollouts"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-rollouts"
  create_namespace = true
  wait             = true
  version          = "2.32.7"

  values = [
    <<-EOF
    controller:
      replicas: 1
      metrics:
        enabled: true
        serviceMonitor:
          enabled: ${var.prometheus_enabled}
          additionalLabels:
            release: prometheus
    dashboard:
      enabled: true
      service:
        type: NodePort
        nodePort: ${var.node_port}
    EOF
  ]
}
