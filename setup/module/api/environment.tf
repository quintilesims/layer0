resource "aws_ecs_cluster" "api" {
  name = "l0-${var.name}-api"
}

resource "aws_security_group" "api_env" {
  name        = "l0-${var.name}-api-env"
  description = "Auto-generated Layer0 Environment Security Group"
  vpc_id      = "${var.vpc_id}"
  tags        = "${var.tags}"

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

data "template_file" "user_data" {
  template = "${file("${path.module}/user_data.sh")}"

  vars {
    cluster_id = "${aws_ecs_cluster.api.id}"
    s3_bucket  = "${aws_s3_bucket.mod.id}"
  }
}

resource "aws_launch_configuration" "api" {
  name_prefix          = "l0-${var.name}-api-"
  image_id             = "${data.aws_ami.linux.id}"
  instance_type        = "t3.medium"
  security_groups      = ["${aws_security_group.api_env.id}"]
  iam_instance_profile = "${aws_iam_instance_profile.ecs.id}"
  user_data            = "${data.template_file.user_data.rendered}"
  key_name             = "${var.ssh_key_pair}"
  ebs_optimized        = true

  root_block_device {
    delete_on_termination = true
    volume_type           = "gp3"
    volume_size           = "8"
    iops                  = 3000
    throughput            = 125
  }
  ebs_block_device {
    delete_on_termination = true
    device_name           ="/dev/xvdcz"
    volume_type           = "gp3"
    volume_size           = "22"
    iops                  = 3000
    throughput            = 125
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "api" {
  name                 = "l0-${var.name}-api"
  launch_configuration = "${aws_launch_configuration.api.name}"
  vpc_zone_identifier  = ["${data.aws_subnet_ids.private.ids}"]
  min_size             = "2"
  desired_capacity     = "2"
  max_size             = "2"

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
