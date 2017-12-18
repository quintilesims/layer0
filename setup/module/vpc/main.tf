data "aws_availability_zones" "available" {}

resource "aws_vpc" "mod" {
  count = "${var.count_hack}"

  cidr_block           = "${var.cidr}"
  enable_dns_hostnames = "true"
  enable_dns_support   = "true"
  tags                 = "${merge(var.tags, map("Name", format("l0-%s", var.name)))}"
}

resource "aws_internet_gateway" "mod" {
  count = "${var.count_hack}"

  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-igw", var.name)))}"
}

resource "aws_route_table" "public" {
  count = "${var.count_hack}"

  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-rt-public", var.name)))}"
}

resource "aws_route" "public_internet_gateway" {
  count = "${var.count_hack}"

  route_table_id         = "${aws_route_table.public.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.mod.id}"
}

resource "aws_route" "private_nat_gateway" {
  count = "${var.count_hack}"

  route_table_id         = "${aws_route_table.private.id}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.natgw.id}"
}

resource "aws_route_table" "private" {
  count = "${var.count_hack}"

  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-rt-private", var.name)))}"
}

resource "aws_subnet" "private" {
  vpc_id            = "${aws_vpc.mod.id}"
  cidr_block        = "${cidrsubnet(aws_vpc.mod.cidr_block, 8, count.index + 1)}" 
  availability_zone = "${element(data.aws_availability_zones.available.names, count.index)}"
  count             = "${length(data.aws_availability_zones.available.names) * var.count_hack}"
  tags              = "${merge(var.tags, map("Tier", "Private"), map("Name", format("l0-%s-subnet-private-%s", var.name, element(data.aws_availability_zones.available.names, count.index))))}"
}

resource "aws_subnet" "public" {
  vpc_id            = "${aws_vpc.mod.id}"
  cidr_block        = "${cidrsubnet(aws_vpc.mod.cidr_block, 8, count.index + 1 + 100)}"
  availability_zone = "${element(data.aws_availability_zones.available.names, count.index)}"
  count             = "${length(data.aws_availability_zones.available.names) * var.count_hack}"
  tags              = "${merge(var.tags, map("Tier", "Public"), map("Name", format("l0-%s-subnet-public-%s", var.name, element(data.aws_availability_zones.available.names, count.index))))}"

  map_public_ip_on_launch = "${var.map_public_ip_on_launch}"
}

resource "aws_eip" "nateip" {
  count = "${var.count_hack}"

  vpc = true
}

resource "aws_nat_gateway" "natgw" {
  count = "${var.count_hack}"

  allocation_id = "${aws_eip.nateip.id}"
  subnet_id     = "${element(aws_subnet.public.*.id, 0)}"

  depends_on = ["aws_internet_gateway.mod"]
}

resource "aws_route_table_association" "private" {
  count          = "${length(data.aws_availability_zones.available.names) * var.count_hack}"
  subnet_id      = "${element(aws_subnet.private.*.id, count.index)}"
  route_table_id = "${aws_route_table.private.id}"
}

resource "aws_route_table_association" "public" {
  count          = "${length(data.aws_availability_zones.available.names) * var.count_hack}"
  subnet_id      = "${element(aws_subnet.public.*.id, count.index)}"
  route_table_id = "${aws_route_table.public.id}"
}
