variable "alertmanager_url" {
  type        = string
  description = "AlertManager URL"
}

variable "alerting_rules" {
  type        = list(string)
  default     = []
  description = "Alerting rules"
}

variable "labels" {
  type        = list(string)
  default     = ["level", "msg", "caller"]
  description = "Labels"
}

variable "prometheus_enabled" {
  type        = bool
  default     = true
  description = "Enable Prometheus"
}
