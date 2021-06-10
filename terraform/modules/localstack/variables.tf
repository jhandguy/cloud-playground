variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
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