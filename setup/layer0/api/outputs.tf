output "load_balancer_url" {
  value = "${aws_elb.api.dns_name}"
}

output "token" {
  value = "${base64encode("${var.username}:${var.password}")}"
}

output "public_subnets" {
  value = "${join(",", data.aws_subnet_ids.public.ids)}"
}

output "private_subnets" {
  value = "${join(",", data.aws_subnet_ids.private.ids)}"
}

output "linux_service_ami" {
  value = "${lookup(var.linux_region_amis, var.region)}"
}

output "windows_service_ami" {
  value = "${lookup(var.windows_region_amis, var.region)}"
}
