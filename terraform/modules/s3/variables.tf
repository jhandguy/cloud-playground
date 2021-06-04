variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
}

variable "image_registry" {
  type        = string
  sensitive   = true
  description = "Image registry"
}