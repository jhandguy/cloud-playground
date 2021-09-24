resource "random_password" "root_token" {
  length  = 32
  special = false
}
