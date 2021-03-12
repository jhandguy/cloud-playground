resource "kubernetes_config_map" "cli" {
  metadata {
    name      = "cli"
    namespace = kubernetes_namespace.cli.metadata.0.name
  }

  data = {
    gateway_url     = var.gateway_url
    pushgateway_url = var.pushgateway_url
  }
}