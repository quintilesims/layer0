resource "aws_elb" "api" {
  name            = "l0-${var.name}-api"
  subnets         = ["${var.public_subnets}"]
  security_groups = ["${aws_security_group.api_env.id}", "${aws_security_group.api_lb.id}"]

  listener {
    instance_port      = 80
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "${aws_iam_server_certificate.api.arn}"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    target              = "HTTP:80/admin/health"
    interval            = 30
  }
}

resource "tls_private_key" "api" {
  algorithm = "RSA"
}

resource "aws_iam_server_certificate" "api" {
  name             = "l0-${var.name}-api"
  certificate_body = "${tls_self_signed_cert.api.cert_pem}"
  private_key      = "${tls_private_key.api.private_key_pem}"
}

resource "tls_self_signed_cert" "api" {
  key_algorithm   = "${tls_private_key.api.algorithm}"
  private_key_pem = "${tls_private_key.api.private_key_pem}"

  subject {
    common_name = "example.com"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}
