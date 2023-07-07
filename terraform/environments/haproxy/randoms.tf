resource "random_pet" "postgres_user" {
  length = 1
}

resource "random_pet" "postgres_database" {
  length = 1
}

resource "random_pet" "mysql_user" {
  length = 1
}

resource "random_pet" "mysql_database" {
  length = 1
}

resource "random_password" "sql_postgres_token" {
  length  = 32
  special = false
}

resource "random_password" "sql_mysql_token" {
  length  = 32
  special = false
}

resource "random_pet" "sql_postgres_host" {
  length    = 3
  separator = "."
}
resource "random_pet" "sql_mysql_host" {
  length    = 3
  separator = "."
}
