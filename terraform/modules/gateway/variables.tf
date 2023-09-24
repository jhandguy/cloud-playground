variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = map(tuple([number, number]))
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
  default     = true
  description = "Enable Prometheus"
}

variable "prometheus_url" {
  type        = string
  default     = ""
  description = "Prometheus URL"
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

variable "argorollouts_enabled" {
  type        = bool
  default     = false
  description = "Enable ArgoRollouts"
}
