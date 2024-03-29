variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "alertmanager_node_port" {
  type        = number
  description = "AlertManager node port"
}

variable "grafana_node_port" {
  type        = number
  description = "Grafana node port"
}

variable "prometheus_node_port" {
  type        = number
  description = "Prometheus node port"
}

variable "grafana_datasources" {
  type        = list(string)
  default     = []
  description = "Grafana datasources"
}

variable "grafana_dashboards" {
  type        = list(string)
  default     = []
  description = "Grafana dashboards"
}

variable "mimir_url" {
  type        = string
  default     = ""
  description = "Mimir URL"
}
