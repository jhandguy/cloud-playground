variable "consul_enabled" {
  type        = bool
  default     = false
  description = "Enable Consul"
}

variable "prometheus_enabled" {
  type        = bool
  default     = true
  description = "Enable Prometheus"
}
