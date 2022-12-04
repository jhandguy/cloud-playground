output "alertmanager_url" {
  value       = "${var.node_ip}:${var.alertmanager_node_port}"
  description = "AlertManager URL"
}

output "alertmanager_cluster_url" {
  value       = "${helm_release.prometheus.name}-kube-prometheus-alertmanager.${helm_release.prometheus.namespace}.svc.cluster.local:9093"
  description = "AlertManager Cluster URL"
}

output "grafana_url" {
  value       = "${var.node_ip}:${var.grafana_node_port}"
  description = "Grafana URL"
}

output "prometheus_url" {
  value       = "${var.node_ip}:${var.prometheus_node_port}"
  description = "Prometheus URL"
}

output "prometheus_cluster_url" {
  value       = "${helm_release.prometheus.name}-kube-prometheus-prometheus.${helm_release.prometheus.namespace}.svc.cluster.local:9090"
  description = "Prometheus URL"
}

output "grafana_admin_password" {
  value       = random_password.admin_password.result
  description = "Grafana admin password"
  sensitive   = true
}
