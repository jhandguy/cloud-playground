variable "aws_region" {
  type        = string
  sensitive   = true
  description = "AWS region"
  default     = ""
}

variable "aws_access_key_id" {
  type        = string
  sensitive   = true
  description = "AWS access key id"
  default     = ""
}

variable "aws_secret_access_key" {
  type        = string
  sensitive   = true
  description = "AWS secret access key"
  default     = ""
}

variable "aws_s3_bucket" {
  type        = string
  description = "AWS S3 bucket"
  default     = ""
}

variable "aws_s3_cluster_endpoint" {
  type        = string
  description = "AWS S3 cluster endpoint"
  default     = ""
}

variable "localstack_enabled" {
  type        = bool
  default     = true
  description = "Enable LocalStack"
}

variable "prometheus_enabled" {
  type        = bool
  default     = true
  description = "Enable Prometheus"
}
