output "environment_ids" {
  value = ["${layer0_environment.te.*.id}"]
}

output "load_balancer_ids" {
  value = ["${layer0_load_balancer.tlb.*.id}"]
}

output "deploy_ids" {
  value = ["${layer0_deploy.td.*.id}"]
}

output "service_ids" {
  value = ["${layer0_service.ts.*.id}"]
}
