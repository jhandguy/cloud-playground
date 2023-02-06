variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = tuple([number, number])
  description = "Node ports"
}

variable "prometheus_enabled" {
  type        = bool
  default     = false
  description = "Enable Prometheus"
}

variable "consul_enabled" {
  type        = bool
  default     = false
  description = "Enable Consul"
}

variable "csi_enabled" {
  type        = bool
  default     = false
  description = "Enable CSI"
}

variable "vault_url" {
  type        = string
  default     = ""
  description = "Vault URL"
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

variable "min_replicas" {
  type        = number
  default     = 1
  description = "Minimum replicas"
}

variable "max_replicas" {
  type        = number
  default     = 1
  description = "Maximum replicas"
}
