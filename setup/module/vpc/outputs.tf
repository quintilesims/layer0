# todo: join hack is a workaround for https://github.com/hashicorp/hil/issues/50
output "vpc_id" {
  value =  var.count_hack == 0 ? "<none>" : join(" ", aws_vpc.mod.*.id) 
}
