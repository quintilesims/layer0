data "aws_subnet_ids" "private" {
  vpc_id = "${var.vpc_id}"

  tags {
    Tier = "Private"
  }
}

data "aws_subnet_ids" "public" {
  vpc_id = "${var.vpc_id}"

  tags {
    Tier = "Public"
  }
}

resource "aws_ecs_cluster" "api" {
  name = "l0-${var.name}-api"
}

data "template_file" "user_data" {
    template = "${file("${path.module}/user_data.sh")}"
    vars {
        cluster_id = "${aws_ecs_cluster.api.id}"
        s3_bucket = "${var.bucket_name}"
    }
}
