output "jenkins" {
  value       = "http://${kubernetes_ingress.jenkins.load_balancer_ingress.0.ip}${local.jenkins_uri_prefix}/"
  description = "Jenkins URL"
}

output "user" {
  value       = data.kubernetes_secret.jenkins.data["jenkins-admin-user"]
  description = "Jenkins User"
}

output "password" {
  value       = data.kubernetes_secret.jenkins.data["jenkins-admin-password"]
  description = "Jenkins Password"
}

output "localstack" {
  value       = "http://${kubernetes_ingress.localstack.load_balancer_ingress.0.ip}/"
  description = "LocalStack URL"
}

output "bucket" {
  value       = aws_s3_bucket.bucket.id
  description = "Bucket ID"
}

output "s3" {
  value       = "http://${kubernetes_ingress.s3.load_balancer_ingress.0.ip}${local.s3_uri_prefix}/"
  description = "S3 URL"
}

output "s3_health" {
  value       = "http://${kubernetes_ingress.s3.load_balancer_ingress.0.ip}${local.s3_uri_prefix}${local.s3_health_path}"
  description = "S3 Health URL"
}