output "alertmanager_url" {
  value       = "${var.node_ip}:${var.alertmanager_node_port}"
  description = "AlertManager URL"
}

output "grafana_url" {
  value       = "${var.node_ip}:${var.grafana_node_port}"
  description = "Grafana URL"
}

output "prometheus_url" {
  value       = "${var.node_ip}:${var.prometheus_node_port}"
  description = "Prometheus URL"
}

output "pushgateway_url" {
  value       = "${var.node_ip}:${var.pushgateway_node_port}"
  description = "PushGateway URL"
}

output "grafana_admin_password" {
  value       = random_password.admin_password.result
  description = "Grafana admin password"
  sensitive   = true
}