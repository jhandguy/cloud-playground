resource "random_password" "user_password" {
  length  = 32
  special = false
}

resource "random_password" "postgres_password" {
  length  = 32
  special = false
}
