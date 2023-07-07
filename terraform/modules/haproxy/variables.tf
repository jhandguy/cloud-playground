variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = tuple([number, number, number])
  description = "Node ports"
}

variable "prometheus_enabled" {
  type        = bool
  default     = false
  description = "Enable Prometheus"
}