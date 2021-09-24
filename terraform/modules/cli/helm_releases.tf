resource "helm_release" "cli" {
  name             = "cli"
  namespace        = "cli"
  chart            = "../cli/helm"
  create_namespace = true
  wait             = true
  wait_for_jobs    = true
  version          = "1.0.0"

  values = [
    <<-EOF
    test:
      rounds: ${var.test_rounds}
    EOF
  ]
}
