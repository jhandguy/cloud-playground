output "cluster_url" {
  value       = "${helm_release.tempo.namespace}.${helm_release.tempo.name}.svc.cluster.local:3100"
  description = "Cluster URL"
}

output "otlp_grpc_url" {
  value       = "${helm_release.tempo.namespace}.${helm_release.tempo.name}.svc.cluster.local:4317"
  description = "OTLP gRPC URL"
}
