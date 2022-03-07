output "cluster_url" {
  value       = "${helm_release.loki.name}.${helm_release.loki.namespace}.svc.cluster.local:3100"
  description = "Cluster URL"
}
