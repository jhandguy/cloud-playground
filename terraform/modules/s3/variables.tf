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
