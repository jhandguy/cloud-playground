output "cluster_url" {
  value       = "${helm_release.loki.namespace}.${helm_release.loki.name}.svc.cluster.local:3100"
  description = "Cluster URL"
}
