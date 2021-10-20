variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
}

variable "alertmanager_node_port" {
  type        = number
  description = "AlertManager node port"
}

variable "alerting_rules" {
  type        = list(string)
  default     = []
  description = "Alerting rules"
}
