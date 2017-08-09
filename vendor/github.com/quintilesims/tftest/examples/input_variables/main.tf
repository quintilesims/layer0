variable "message_one" {}

variable "message_two" {}

output "combined_messages" {
  value = "${var.message_one} ${var.message_two}"
}
