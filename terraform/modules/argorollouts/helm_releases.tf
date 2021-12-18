resource "helm_release" "argorollouts" {
  name             = "argorollouts"
  namespace        = "argorollouts"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-rollouts"
  create_namespace = true
  wait             = true
  version          = "2.6.0"

  values = [
    <<-EOF
    controller:
%{if var.prometheus_enabled}
      metrics:
        enabled: true
        serviceMonitor:
          enabled: true
          additionalLabels:
            release: prometheus
%{endif}
    dashboard:
      enabled: true
    EOF
  ]
}