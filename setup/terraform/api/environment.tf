resource "aws_security_group" "api-env" {
  name        = "l0-${var.layer0_instance_name}-api-env"
  description = "Auto-generated Layer0 Environment Security Group"
  vpc_id      = "${var.vpc_id}"

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

  tags {
    "layer0" = "${var.layer0_instance_name}"
  }
}

resource "aws_ecs_cluster" "api" {
  name = "l0-${var.layer0_instance_name}-api"
}

resource "aws_launch_configuration" "api" {
  name_prefix          = "l0-${var.layer0_instance_name}-api-"
  image_id             = "${var.ami_id}"
  instance_type        = "${var.instance_type}"
  iam_instance_profile = "${aws_iam_instance_profile.api.id}"
  security_groups      = ["${aws_security_group.api-env.id}"]
  user_data            = "${data.template_file.user_data.rendered}"
  key_name             = "${var.key_pair}"

  tags {
    "layer0" = "${var.layer0_instance_name}"
  }

  # see: https://www.terraform.io/docs/configuration/resources.html
  lifecycle {
    create_before_destroy = true
  }
}

data "template_file" "user_data" {
  template = "${file("${path.module}/user_data.sh")}"

  vars {
    cluster_id = "${aws_ecs_cluster.api.id}"
    s3_bucket  = "${var.s3_bucket}"
  }
}

resource "aws_autoscaling_group" "api" {
  name                      = "l0-${var.layer0_instance_name}-api"
  launch_configuration      = "${aws_launch_configuration.api.name}"
  vpc_zone_identifier       = ["${aws_subnet.private_secondary.id}"]
  min_size                  = "2"
  max_size                  = "2"
  desired_capacity          = "2"
  health_check_type         = "EC2"
  health_check_grace_period = "300"

  depends_on = [
    "aws_db_instance.api",
  ]

  tag {
    key                 = "Name"
    value               = "l0-${var.layer0_instance_name}-api"
    propagate_at_launch = true
  }

  tags {
    "layer0" = "${var.layer0_instance_name}"
  }

  lifecycle {
    create_before_destroy = true
  }
}
