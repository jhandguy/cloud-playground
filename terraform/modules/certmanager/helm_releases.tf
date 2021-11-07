resource "helm_release" "certmanager" {
  name             = "certmanager"
  namespace        = "certmanager"
  repository       = "https://charts.jetstack.io"
  chart            = "cert-manager"
  create_namespace = true
  wait             = true
  version          = "1.5.4"

  values = [
    <<-EOF
    installCRDs: true
    startupapicheck:
      enabled: false
    EOF
  ]
}
