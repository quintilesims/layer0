output "load_balancer_url" {
  value = "${aws_elb.api.dns_name}"
}
