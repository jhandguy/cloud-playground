variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = map(number)
  default     = {}
  description = "Node ports"
}

variable "ingress_gateway_port" {
  type        = number
  default     = 8080
  description = "Ingress Gateway port"
}

variable "ingress_host" {
  type        = string
  description = "Ingress host"
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
