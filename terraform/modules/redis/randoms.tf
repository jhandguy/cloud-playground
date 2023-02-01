resource "random_password" "redis_password" {
  length  = 32
  special = false
}
