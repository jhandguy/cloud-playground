resource "random_shuffle" "node_ports" {
  result_count = length(var.node_ports)
  input        = range(30000, 30000 + length(var.node_ports))
}
