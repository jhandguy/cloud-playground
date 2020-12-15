output "jenkins" {
  value       = "http://${kubernetes_ingress.jenkins.load_balancer_ingress.0.ip}/${helm_release.jenkins.name}/"
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