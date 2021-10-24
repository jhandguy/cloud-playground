resource "helm_release" "metrics" {
  name             = "metrics"
  namespace        = "metrics"
  repository       = "https://kubernetes-sigs.github.io/metrics-server"
  chart            = "metrics-server"
  create_namespace = true
  wait             = true
  version          = "3.6.0"

  values = [
    <<-EOF
    args:
      - --metric-resolution=10s
      - --kubelet-insecure-tls
    EOF
  ]
}
