output "token" {
  value       = random_password.dynamo_token.result
  sensitive   = true
  description = "Token"
}