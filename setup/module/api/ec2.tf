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
