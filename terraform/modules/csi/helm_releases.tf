resource "helm_release" "csi" {
  name             = "csi"
  namespace        = "csi"
  repository       = "https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts"
  chart            = "secrets-store-csi-driver"
  create_namespace = true
  wait             = true
  version          = "0.2.0"

  values = [
    <<-EOF
    syncSecret:
      enabled: true
    EOF
  ]
}
