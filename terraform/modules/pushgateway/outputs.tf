output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.pushgateway.name}-prometheus-pushgateway.${helm_release.pushgateway.namespace}.svc.cluster.local:9091"
  description = "Cluster URL"
}
