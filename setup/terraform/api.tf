module "api" {
  source = "${path.module}/api"
  name = "l0-${layer0_instance_name}"
}

