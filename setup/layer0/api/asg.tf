resource "aws_launch_configuration" "api" {
  name_prefix          = "l0-${var.name}-api-"
  image_id             = "${var.ami}"
  instance_type        = "t2.small"
  iam_instance_profile = "${var.instance_profile}"
  security_groups      = ["${aws_security_group.api_env.id}"]
  user_data            = "${data.template_file.user_data.rendered}"
  key_name             = "${var.key_pair}"

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

resource "aws_autoscaling_group" "ecs_cluster" {
  name                      = "l0-${var.name}-api"
  launch_configuration      = "${aws_launch_configuration.api.name}"
  vpc_zone_identifier       = ["${var.public_subnets}"]
  min_size                  = 2
  max_size                  = 2
  desired_capacity          = 2
  health_check_type         = "EC2"
  health_check_grace_period = "300"

  tag {
    key                 = "Name"
    value               = "l0-${var.name}-api"
    propagate_at_launch = true
  }

  tag {
    key                 = "layer0"
    value               = "${var.name}"
    propagate_at_launch = true
  }

  lifecycle {
    create_before_destroy = true
  }
}
