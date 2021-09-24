output "urls" {
  value = {
    for name, node_port in var.node_ports : name => "${var.node_ip}:${node_port}"
  }
  description = "URLs"
}
