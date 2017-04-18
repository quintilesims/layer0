resource "aws_ecs_cluster" "api" {
  name = "l0-${var.name}-api"
}

resource "aws_ecs_service" "api" {
  name            = "l0-${var.name}-api"
  cluster         = "${aws_ecs_cluster.api.id}"
  task_definition = "${aws_ecs_task_definition.api.arn}"
  desired_count   = 1
  iam_role        = "${var.ecs_role_arn}"
  # TODO: docs have: depends_on      = ["aws_iam_role_policy.ecs"]

  deployment_minimum_healthy_percent = 0
  deployment_minimum_healthy_percent = 200

  load_balancer {
    elb_name       = "${aws_elb.api.name}"
    container_name = "api"
    container_port = 9090
  }
}

resource "aws_ecs_task_definition" "api" {
  family                = "l0-${var.name}-api"
  container_definitions = "${data.template_file.container_definitions.rendered}"
}

data "template_file" "container_definitions" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    name           = "${var.name}"
    region         = "${var.aws_region}"
    api_image_tag  = "todo"
    log_group_name = "todo"
  }
}

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

resource "aws_security_group" "api_env" {
  name        = "l0-${var.name}-api-env"
  description = "Auto-generated Layer0 Environment Security Group"
  vpc_id      = "${var.vpc}"

  ingress {
    self      = "true"
    from_port = 0
    to_port   = 0
    protocol  = "-1"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "api_lb" {
  name        = "l0-${var.name}-api-lb"
  description = "Auto-generated Layer0 Load Balancer Security Group"
  vpc_id      = "${var.vpc}"

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
