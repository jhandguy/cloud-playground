resource "grafana_data_source" "loki" {
  depends_on = [helm_release.prometheus]

  type = "loki"
  name = "Loki"
  url  = "http://${var.node_ip}:${var.loki_node_port}"
}