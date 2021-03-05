provider "grafana" {
  url  = "http://${var.node_ip}:${var.grafana_node_port}"
  auth = "admin:${random_password.admin_password.result}"
}