output "name" {
  value = "${var.name}"
}

output "s3_bucket" {
  value = "todo"
}

output "vpc_id" {
  value = "${var.vpc_id != "" ? var.vpc_id : "the vpc we created!"}"
}
