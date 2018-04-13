variable "environment_id" {}

variable "private" {
  default = false
}

variable "scale" {
  default = 1
}

variable "load_balancer_type" {
    default = "application"
}

variable "stateful" {
  default = false
}
