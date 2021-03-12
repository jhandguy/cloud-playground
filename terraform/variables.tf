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

variable "s3_image_repository" {
  type        = string
  sensitive   = true
  description = "S3 image repository"
}

variable "s3_image_tag" {
  type        = string
  sensitive   = true
  description = "S3 image tag"
}

variable "dynamo_image_repository" {
  type        = string
  sensitive   = true
  description = "Dynamo image repository"
}

variable "dynamo_image_tag" {
  type        = string
  sensitive   = true
  description = "Dynamo image tag"
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

variable "cli_image_repository" {
  type        = string
  sensitive   = true
  description = "CLI image repository"
}

variable "cli_image_tag" {
  type        = string
  sensitive   = true
  description = "CLI image tag"
}