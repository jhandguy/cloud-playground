variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
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
