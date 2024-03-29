output "url" {
  value       = "${var.node_ip}:${var.node_ports.0}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.s3.name}.${helm_release.s3.namespace}.svc.cluster.local:8080"
  description = "Cluster URL"
}
