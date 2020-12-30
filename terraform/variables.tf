variable "aws_region" {
  type      = string
  sensitive = true
}

variable "aws_access_key_id" {
  type      = string
  sensitive = true
}

variable "aws_secret_access_key" {
  type      = string
  sensitive = true
}

variable "registry_username" {
  type      = string
  sensitive = true
}

variable "registry_password" {
  type      = string
  sensitive = true
}

variable "image_registry" {
  type      = string
  sensitive = true
}

variable "image_repository" {
  type      = string
  sensitive = true
}

variable "image_tag" {
  type      = string
  sensitive = true
}