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

variable "argorollouts_enabled" {
  type        = bool
  default     = false
  description = "Enable ArgoRollouts"
}
