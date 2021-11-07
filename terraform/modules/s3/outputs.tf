output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.s3.namespace}.${helm_release.s3.name}.svc.cluster.local:8080"
  description = "Cluster URL"
}
