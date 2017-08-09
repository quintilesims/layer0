variable "message" {
  default = "Hello World"
}

output "message" {
  value = "${var.message}"
}
