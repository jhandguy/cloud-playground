variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
}

variable "prometheus_enabled" {
  type        = bool
  default     = false
  description = "Enable Prometheus"
}
