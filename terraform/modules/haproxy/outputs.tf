output "url" {
  value       = "${var.node_ip}:${var.node_ports.0}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.haproxy.name}-kubernetes-ingress.${helm_release.haproxy.namespace}.svc.cluster.local:80"
  description = "Cluster URL"
}
