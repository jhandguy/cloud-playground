output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.mimir.name}-gateway.${helm_release.mimir.namespace}.svc.cluster.local:80"
  description = "Cluster URL"
}
