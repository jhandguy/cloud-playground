output "cluster_context" {
  value       = "kind-${var.cluster_name}"
  description = "Cluster context"
}

output "node_ip" {
  value       = "localhost"
  description = "Node ip"
}

output "node_ports" {
  value = {
    for index in range(length(var.node_ports)) : var.node_ports[index] => random_shuffle.node_ports.result[index]
  }
  description = "Node ports"
}
