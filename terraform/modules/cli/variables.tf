variable "test_rounds" {
  type        = number
  description = "Test rounds"
  default     = 100
}

variable "csi_enabled" {
  type        = bool
  default     = false
  description = "Enable CSI"
}

variable "vault_url" {
  type        = string
  default     = ""
  description = "Vault URL"
}

variable "secrets" {
  type        = map(string)
  default     = {}
  description = "Secrets"
}
