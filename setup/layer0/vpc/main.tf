resource "aws_vpc" "mod" {
  cidr_block           = "${var.cidr}"
  enable_dns_hostnames = "true"
  enable_dns_support   = "true"
  tags                 = "${merge(var.tags, map("Name", format("l0-%s", var.name)))}"
}

resource "aws_internet_gateway" "mod" {
  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-igw", var.name)))}"
}

resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-rt-public", var.name)))}"
}

resource "aws_route" "public_internet_gateway" {
  route_table_id         = "${aws_route_table.public.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.mod.id}"
}

resource "aws_route" "private_nat_gateway" {
  route_table_id         = "${aws_route_table.private.id}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.natgw.id}"
}

resource "aws_route_table" "private" {
  vpc_id = "${aws_vpc.mod.id}"
  tags   = "${merge(var.tags, map("Name", format("l0-%s-rt-private", var.name)))}"
}

resource "aws_subnet" "private" {
  vpc_id            = "${aws_vpc.mod.id}"
  cidr_block        = "${var.private_subnets[count.index]}"
  availability_zone = "${element(var.azs, count.index)}"
  count             = "${length(var.private_subnets)}"
  tags              = "${merge(var.tags, map("Tier", "Private"), map("Name", format("l0-%s-subnet-private-%s", var.name, element(var.azs, count.index))))}"
}

resource "aws_subnet" "public" {
  vpc_id            = "${aws_vpc.mod.id}"
  cidr_block        = "${var.public_subnets[count.index]}"
  availability_zone = "${element(var.azs, count.index)}"
  count             = "${length(var.public_subnets)}"
  tags              = "${merge(var.tags, map("Tier", "Public"), map("Name", format("l0-%s-subnet-public-%s", var.name, element(var.azs, count.index))))}"

  map_public_ip_on_launch = "${var.map_public_ip_on_launch}"
}

resource "aws_eip" "nateip" {
  vpc = true
}

resource "aws_nat_gateway" "natgw" {
  allocation_id = "${aws_eip.nateip.id}"
  subnet_id     = "${element(aws_subnet.public.*.id, 0)}"

  depends_on = ["aws_internet_gateway.mod"]
}

resource "aws_route_table_association" "private" {
  count          = "${length(var.private_subnets)}"
  subnet_id      = "${element(aws_subnet.private.*.id, count.index)}"
  route_table_id = "${aws_route_table.private.id}"
}

resource "aws_route_table_association" "public" {
  count          = "${length(var.public_subnets)}"
  subnet_id      = "${element(aws_subnet.public.*.id, count.index)}"
  route_table_id = "${aws_route_table.public.id}"
}
