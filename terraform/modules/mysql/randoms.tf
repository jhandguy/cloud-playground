resource "random_password" "user_password" {
  length  = 32
  special = false
}

resource "random_password" "root_password" {
  length  = 32
  special = false
}
