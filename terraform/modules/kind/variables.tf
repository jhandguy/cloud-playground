variable "cluster_name" {
  type        = string
  description = "Cluster name"
  default     = "kind"
}

variable "node_image" {
  type        = string
  description = "Node image"
  default     = "v1.25.3"
}

variable "node_ports" {
  type        = list(string)
  description = "Node ports"
  default     = []
}
