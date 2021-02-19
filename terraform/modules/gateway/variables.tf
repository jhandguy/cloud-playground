variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
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

variable "gateway_image_repository" {
  type        = string
  sensitive   = true
  description = "Gateway image repository"
}

variable "gateway_image_tag" {
  type        = string
  sensitive   = true
  description = "Gateway image tag"
}

variable "s3_url" {
  type        = string
  description = "S3 URL"
}

variable "s3_token" {
  type        = string
  sensitive   = true
  description = "S3 token"
}

variable "dynamo_url" {
  type        = string
  description = "Dynamo URL"
}

variable "dynamo_token" {
  type        = string
  sensitive   = true
  description = "Dynamo token"
}