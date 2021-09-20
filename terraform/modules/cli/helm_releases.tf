resource "helm_release" "cli" {
  name             = "cli"
  namespace        = "cli"
  chart            = "../cli/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [<<-EOF
    test:
      rounds: ${var.test_rounds}
    EOF
  ]

  provisioner "local-exec" {
    command = <<-EOF
      kubectl wait --for=condition=complete --timeout=60s job/cli -n cli
      kubectl logs job/cli -c cli -n cli
    EOF
  }
}