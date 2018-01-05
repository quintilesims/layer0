variable "count_hack" {}

variable "name" {}

variable "cidr" {}

variable "map_public_ip_on_launch" {
  description = "Auto-assign public IP on launch"
  default     = false
}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}
