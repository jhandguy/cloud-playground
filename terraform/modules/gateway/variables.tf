variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_ports" {
  type        = map(number)
  description = "Node ports"
}

variable "ingress_gateway_port" {
  type        = number
  description = "Ingress Gateway port"
}

variable "image_registry" {
  type        = string
  sensitive   = true
  description = "Image registry"
}