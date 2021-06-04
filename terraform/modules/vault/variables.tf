variable "node_ip" {
  type        = string
  description = "Node ip"
}

variable "node_port" {
  type        = number
  description = "Node port"
}

variable "secrets" {
  type        = map(map(string))
  description = "Secrets"
}