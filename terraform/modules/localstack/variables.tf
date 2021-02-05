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
  description = "AWS secret access key"
}

variable "aws_s3_buckets" {
  type        = list(string)
  default     = []
  description = "AWS S3 buckets"
}

variable "aws_dynamo_tables" {
  type        = list(string)
  default     = []
  description = "AWS DynamoDB tables"
}