resource "helm_release" "cli" {
  name             = "cli"
  namespace        = "cli"
  chart            = "../../../cli/helm"
  create_namespace = true
  wait             = true
  wait_for_jobs    = true
  version          = "1.0.0"

  values = [
    <<-EOF
    test:
      rounds: ${var.test_rounds}
    csi:
      enabled: ${var.csi_enabled}
      vaultAddress: ${var.vault_url}
    EOF
  ]

  dynamic "set_sensitive" {
    for_each = var.secrets

    content {
      name  = "secrets.${set_sensitive.key}"
      value = base64encode(set_sensitive.value)
    }
  }
}
