resource "helm_release" "cli" {
  name      = "cli"
  namespace = kubernetes_namespace.cli.metadata.0.name
  chart     = "../cli/helm"
  wait      = true

  values = [<<-EOF
    configMap: ${kubernetes_config_map.cli.metadata.0.name}
    secret: ${kubernetes_secret.cli.metadata.0.name}
    image:
      secret: ${kubernetes_secret.cli_image.metadata.0.name}
      registry: ${var.image_registry}
    test:
      rounds: 50
    EOF
  ]
}