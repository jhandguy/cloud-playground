output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.nginx.name}-ingress-nginx-controller.${helm_release.nginx.namespace}.svc.cluster.local:80"
  description = "Cluster URL"
}
