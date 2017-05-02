output "name" {
  value = "${var.name}"
}

output "endpoint" {
  value = "https://${module.api.load_balancer_url}"
}

output "s3_bucket" {
  value = "${module.core.bucket_name}"
}

