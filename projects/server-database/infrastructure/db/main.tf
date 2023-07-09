terraform {
  required_providers {
    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = "1.19.0"
    }
  }
}

provider "postgresql" {
  host            = var.host
  port            = var.port
  username        = var.username
  password        = var.password
  connect_timeout = var.connection_timeout
  sslmode         = var.ssl_mode
}

resource "postgresql_database" "go_server_database" {
  provider          = "postgresql"
  name              = "go-server-database"
  owner             = var.username
  lc_collate        = "en_US.UTF-8"
  lc_ctype          = "en_US.UTF-8"
  connection_limit  = -1
  allow_connections = true
  is_template       = false
}