variable "host" {
  type = string
}

variable "port" {
  type    = number
  default = 5432
}

variable "username" {
  type    = string
  default = "root"
}

variable "password" {
  type      = string
  sensitive = true
}

variable "connection_timeout" {
  type    = number
  default = 15
}

variable "ssl_mode" {
  type    = string
  default = "disable"
}