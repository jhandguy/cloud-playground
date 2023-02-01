output "url" {
  value       = "${var.node_ip}:${var.node_port}"
  description = "URL"
}

output "cluster_url" {
  value       = "${helm_release.redis.name}-master.${helm_release.redis.namespace}.svc.cluster.local:6379"
  description = "Cluster URL"
}

output "redis_password" {
  value       = random_password.redis_password.result
  description = "Redis password"
}
