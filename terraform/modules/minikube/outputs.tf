output "node_ports" {
  value = {
    for index in range(0, length(var.node_ports)) : var.node_ports[index] => random_shuffle.node_ports.result[index]
  }
  description = "Node ports"
}