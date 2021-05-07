variable "registry_username" {
  type        = string
  sensitive   = true
  description = "Registry username"
}

variable "registry_password" {
  type        = string
  sensitive   = true
  description = "Registry password"
}

variable "image_registry" {
  type        = string
  sensitive   = true
  description = "Image registry"
}

variable "gateway_api_key" {
  type        = string
  sensitive   = true
  description = "Gateway API key"
}

variable "gateway_url" {
  type        = string
  description = "Gateway URL"
}

variable "pushgateway_url" {
  type        = string
  description = "PushGateway URL"
}