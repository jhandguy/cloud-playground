variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = tuple([number, number])
  description = "Node ports"
}

variable "feature" {
  type        = string
  description = "Feature"

  validation {
    condition     = contains(["postgres", "mysql"], var.feature)
    error_message = "Feature must be one of: [postgres, mysql]"
  }
}

variable "secrets" {
  type        = map(string)
  default     = {}
  description = "Secrets"
}

variable "replicas" {
  type        = number
  default     = 1
  description = "Replicas"
}

variable "prometheus_enabled" {
  type        = bool
  default     = false
  description = "Enable Prometheus"
}
