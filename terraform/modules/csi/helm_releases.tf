resource "helm_release" "csi" {
  name             = "csi"
  namespace        = "csi"
  repository       = "https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts"
  chart            = "secrets-store-csi-driver"
  create_namespace = true
  wait             = true
  version          = "1.2.3"

  values = [
    <<-EOF
    linux:
      driver:
        resources: null
      registrar:
        resources: null
      livenessProbe:
        resources: null
    syncSecret:
      enabled: true
    EOF
  ]
}
