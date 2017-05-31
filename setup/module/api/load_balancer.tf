resource "aws_security_group" "api_lb" {
  name        = "l0-${var.name}-api-lb"
  description = "Auto-generated Layer0 Load Balancer Security Group"
  vpc_id      = "${var.vpc_id}"
  tags        = "${var.tags}"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "tls_private_key" "api" {
  algorithm = "RSA"
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

resource "aws_iam_server_certificate" "api" {
  name             = "l0-${var.name}-api"
  path             = "/l0/l0-${var.name}/"
  certificate_body = "${tls_self_signed_cert.api.cert_pem}"
  private_key      = "${tls_private_key.api.private_key_pem}"
}

resource "aws_elb" "api" {
  name            = "l0-${var.name}-api"
  subnets         = ["${data.aws_subnet_ids.public.ids}"]
  security_groups = ["${aws_security_group.api_env.id}", "${aws_security_group.api_lb.id}"]
  tags            = "${var.tags}"

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
    interval            = 6
  }
}
