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

variable "s3_token" {
  type        = string
  sensitive   = true
  description = "S3 token"
}

variable "dynamo_token" {
  type        = string
  sensitive   = true
  description = "Dynamo token"
}