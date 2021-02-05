variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
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
  description = "AWS secret accest key"
}

variable "aws_s3_endpoint" {
  type        = string
  description = "AWS s3 endpoint"
}

variable "aws_s3_bucket" {
  type        = string
  description = "AWS s3 bucket"
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