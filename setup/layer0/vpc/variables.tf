variable "name" {}

variable "cidr" {}

variable "public_subnets" {
  description = "A list of public subnets inside the VPC."
  default     = []
}

variable "private_subnets" {
  description = "A list of private subnets inside the VPC."
  default     = []
}

variable "azs" {
  description = "A list of Availability zones in the region"
  default     = []
}

variable "map_public_ip_on_launch" {
  description = "Auto-assign public IP on launch"
  default     = false
}

variable "tags" {
  description = "A map of tags to add to all resources"
  default     = {}
}
