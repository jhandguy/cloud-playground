output "token" {
  value       = random_password.s3_token.result
  sensitive   = true
  description = "Token"
}