resource "random_password" "admin_password" {
  length  = 32
  special = false
}
