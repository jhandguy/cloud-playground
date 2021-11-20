variable "alertmanager_url" {
  type        = string
  description = "AlertManager URL"
}

variable "alerting_rules" {
  type        = list(string)
  default     = []
  description = "Alerting rules"
}
