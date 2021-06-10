variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "aws_region" {
  type        = string
  sensitive   = true
  description = "AWS region"
}

variable "aws_access_key_id" {
  type        = string
  sensitive   = true
  description = "AWS access key id"
}

variable "aws_secret_access_key" {
  type        = string
  sensitive   = true
  description = "AWS secret access key"
}

variable "image_registry" {
  type        = string
  sensitive   = true
  description = "Image registry"
  default     = "ghcr.io"
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