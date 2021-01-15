output "localstack" {
  value       = "http://${kubernetes_ingress.localstack.load_balancer_ingress.0.ip}"
  description = "LocalStack URL"
}

output "bucket" {
  value       = aws_s3_bucket.bucket.id
  description = "Bucket ID"
}

output "table" {
  value       = aws_dynamodb_table.table.id
  description = "Table ID"
}

output "s3" {
  value       = "http://${kubernetes_ingress.s3.load_balancer_ingress.0.ip}/${random_pet.s3_uri_prefix.id}"
  description = "S3 URL"
}

output "s3_health" {
  value       = "http://${kubernetes_ingress.s3.load_balancer_ingress.0.ip}/${random_pet.s3_uri_prefix.id}/${random_pet.s3_health_path.id}"
  description = "S3 Health URL"
}

output "dynamo" {
  value       = "http://${kubernetes_ingress.dynamo.load_balancer_ingress.0.ip}/${random_pet.dynamo_uri_prefix.id}"
  description = "Dynamo URL"
}

output "dynamo_health" {
  value       = "http://${kubernetes_ingress.dynamo.load_balancer_ingress.0.ip}/${random_pet.dynamo_uri_prefix.id}/${random_pet.dynamo_health_path.id}"
  description = "Dynamo Health URL"
}