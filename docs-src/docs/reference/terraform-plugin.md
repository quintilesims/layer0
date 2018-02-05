# Layer0 Terraform Provider Reference

Terraform is an open-source tool for provisioning and managing infrastructure.
If you are new to Terraform, we recommend checking out their [documentation](https://www.terraform.io/intro/index.html).

Layer0 has built a custom [provider](https://www.terraform.io/docs/providers/index.html) for Layer0.
This provider allows users to create, manage, and update Layer0 entities using Terraform.

## Prerequisites
- **Terraform v0.11+** ([download](https://www.terraform.io/downloads.html)), accessible in your system path.

## Install
Download a Layer0 v0.8.4+ [release](/releases).
The Terraform plugin binary is located in the release zip file as `terraform-provider-layer0`.

* On Windows, copy `terraform-provider-layer0` in the sub-path `terraform.d\plugins\` beneath your user's `AppData` directory (if running Windows 7 or later). An alias for this path is `%USERPROFILE%\AppData\terraform.d\plugins\`
* On all other systems, copy `terraform-provider-layer0` to the path `~/terraform.d/plugins/` within your user's home directory.

If the `plugins` directory does not exist, create it.

For further information, see Terraform's documentation on installing a Terraform plugin [here](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) and information on third-party plugins [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)

## Getting Started

* Check out the **Terraform** section of the Guestbook walkthrough [here](../guides/walkthrough/deployment-1/#deploy-with-terraform).
* We've added some tips and links to helpful resources in the [Best Practices](#best-practices) section below.

---

## Provider

The Layer0 provider is used to interact with a Layer0 API.
The provider needs to be configured with the proper credentials before it can be used.

### Example Usage

```
# Add 'endpoint' and 'token' variables
variable "endpoint" {}

variable "token" {}

# Configure the layer0 provider
provider "layer0" {
  endpoint        = "${var.endpoint}"
  token           = "${var.token}"
  skip_ssl_verify = true
}
```

### Argument Reference
The following arguments are supported:

!!! note
	The `endpoint` and `token` variables for your layer0 api can be found using the [l0-setup endpoint](setup-cli/#endpoint) command.

* `endpoint` - (Required) The endpoint of the layer0 api
* `token` - (Required) The authentication token for the layer0 api
* `skip_ssl_verify` - (Optional) If true, ssl certificate mismatch warnings will be ignored

---

##API Data Source
The API data source is used to extract useful read-only variables from the Layer0 API.

### Example Usage

```
# Configure the api data source
data "layer0_api" "config" {}

# Output the layer0 vpc id
output "vpc id" {
  val = "${data.layer0_api.config.vpc_id}"
}
```

### Attribute Reference

The following attributes are exported:

* `prefix` - The prefix of the layer0 instance
* `vpc_id` - The vpc id of the layer0 instance
* `public_subnets` - A list containing the 2 public subnet ids in the layer0 vpc
* `private_subnets` - A list containing the 2 private subnet ids in the layer0 vpc

---

##Deploy Data Source

The Deploy data source is used to extract Layer0 Deploy attributes.

### Example Usage

```
# Configure the deploy data source
data "layer0_deploy" "dpl" {
  name    = "my-deploy"
  version = "1"
}

# Output the layer0 deploy id
output "deploy_id" {
  val = "${data.layer0_deploy.dpl.id}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the deploy
* `version` - (Required) The version of the deploy

### Attribute Reference

The following attributes are exported:

* `name` - The name of the deploy
* `version` - The version of the deploy
* `id` - The id of the deploy

---

## Environment Data Source

The Environment data source is used to extract Layer0 Environment attributes.

### Example Usage

```
# Configure the environment data source
data "layer0_environment" "env" {
  name = "my-environment"
}

# Output the layer0 environment id
output "environment_id" {
  val = "${data.layer0_environment.env.id}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment

### Attribute Reference

The following attributes are exported:

* `id` - The id of the environment
* `name` - The name of the environment
* `environment_type` - Default: dynamic. Can be either static or dynamic.
* `instance_type` - The ec2 instance size of a static environment
* `scale` - The current number instances in a static type environment
* `os` - The operating system used for the environment
* `ami` - The AMI ID used for a static environment

---

## Load Balancer Data Source

The Load Balancer data source is used to extract Layer0 Load Balancer attributes.

### Example Usage

```
# Configure the load balancer source
data "layer0_load_balancer" "lb" {
  name           = "my-loadbalancer"
  environment_id = "${data.layer0_environment.env.environment_id}"
}

# Output the layer0 load balancer id
output "load_balancer_id" {
  val = "${data.layer0_load_balancer.lb.id}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (required) The name of the load balancer
* `environment_id` - (required) The id of the environment the load balancer exists in

### Attribute Reference

The following attributes are exported:

* `id` - The id of the load balancer
* `name` - The name of the load balancer
* `environment_id` - The id of the environment the load balancer exists in
* `environment_name` - The name of the environment the load balancer exists in
* `private` - Whether or not the load balancer is private
* `url` - The URL of the load balancer

---

## Service Data Source

The Service data source is used to extract Layer0 Service attributes.

### Example Usage

```
# Configure the service data source
data "layer0_service" "svc" {
  name           = "my-service"
  environment_id = "${data.layer0_environment.env.environment_id}"
}

# Output the layer0 service id
output "service_id" {
  val = "${data.layer0_service.svc.id}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (required) The name of the service
* `environment_id` - (required) The id of the environment the service exists in

### Attribute Reference

The following attributes are exported:

* `id` - The id of the service
* `name` - The name of the service
* `environment_id` - The id of the environment the service exists in
* `environment_name` - The name of the environment the service exists in
* `scale` - The current desired scale of the service

---

## Deploy Resource

Provides a Layer0 Deploy.

Performing variable substitution inside of your deploy's json file (typically named `Dockerrun.aws.json`) can be done through Terraform's [template_file](https://www.terraform.io/docs/providers/template/).
For a working example, please see the sample [Guestbook](https://github.com/quintilesims/guides/blob/master/guestbook/module/main.tf) application

### Example Usage

```
# Configure the deploy template
data "template_file" "guestbook" {
  template = "${file("Dockerrun.aws.json")}"
  vars {
    docker_image_tag = "latest"
  }
}

# Create a deploy using the rendered template
resource "layer0_deploy" "guestbook" {
  name    = "guestbook"
  content = "${data.template_file.guestbook.rendered}"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the deploy
* `content` - (Required) The content of the deploy

### Attribute Reference

The following attributes are exported:

* `id` - The id of the deploy
* `name` - The name of the deploy
* `version` - The version number of the deploy

---

## Environment Resource

Provides a Layer0 Environment

### Example Usage

```
# Create a new environment
resource "layer0_environment" "demo" {
  name             = "demo"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment
* `environment_type` - (Optional, Values: "dynamic, static" Default: "dynamic") Use static when you need to customize ec2 instances that the containers run on or for a Windows os environment.
* `instance_type` - (Optional, Default: "m3.medium") The size of the instances in the environment.
Available instance sizes can be found [here](https://aws.amazon.com/ec2/instance-types/)
* `scale` - (Optional, Default: 0) The desired number of instances allowed in a static environment.
* `user-data` - (Optional) The user data template to use for a static environment's autoscaling group.
See the [cli reference](/reference/cli/#environment) for the default template.
* `os` - (Optional, Default: "linux") Specifies the type of operating system used in the environment.
Options are "linux" or "windows".
* `ami` - (Optional) A custom AMI ID to use for a static environment. 
If not specified, Layer0 will use its default AMI ID for the specified operating system.

### Attribute Reference

The following attributes are exported:

* `id` - The id of the environment
* `name` - The name of the environment
* `environment_type` - The current number instances in a static type environment
* `instance_type` - The ec2 instance size of a static environment
* `scale` - The current number instances in a static type environment
* `security_group_id` - The ID of the environment's security group
* `os` - The operating system used for the environment
* `ami` - The AMI ID used for the environment
---

## Load Balancer Resource

Provides a Layer0 Load Balancer

### Example Usage

```
# Create a new load balancer
resource "layer0_load_balancer" "guestbook" {
  name        = "guestbook"
  environment = "demo123"
  private     = false

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }

  port {
    host_port      = 443
    container_port = 443
    protocol       = "https"
    certificate    = "cert"
  }

  health_check {
    target              = "tcp:80"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the load balancer
* `environment` - (Required) The id of the environment to place the load balancer inside of
* `private` - (Optional) If true, the load balancer will not be exposed to the public internet
* `port` - (Optional, Default: 80:80/tcp) A list of port blocks. Ports documented below
* `health_check` - (Optional, Default: `{"TCP:80" 30 5 2 2}`) A health_check block. Health check documented below

Ports (`port`) support the following:

* `host_port` - (Required) The port on the load balancer to listen on
* `container_port` - (Required) The port on the docker container to route to
* `protocol` - (Required) The protocol to listen on. Valid values are `HTTP, HTTPS, TCP, or SSL`
* `certificate` - (Optional) The name of an SSL certificate. Only required if the `HTTP` or `SSL` protocol is used.

Healthcheck (`health_check`) supports the following:

* `target` - (Required) The target of the check. Valid pattern is `PROTOCOL:PORT/PATH`, where `PROTOCOL` values are:
    * `HTTP`, `HTTPS` - `PORT` and `PATH` are required
    * `TCP`, `SSL` - `PORT` is required, `PATH` is not supported
* `interval` - (Required) The interval between checks.
* `timeout` - (Required) The length of time before the check times out.
* `healthy_threshold` - (Required) The number of checks before the instance is declared healthy.
* `unhealthy_threshold` - (Required) The number of checks before the instance is declared unhealthy.

### Attribute Reference

The following attributes are exported:

* `id` - The id of the load balancer
* `name` - The name of the load balancer
* `environment` - The id of the environment the load balancer exists in
* `private` - Whether or not the load balancer is private
* `url` - The URL of the load balancer

---

## Service Resource

Provides a Layer0 Service

### Example Usage

```
# Create a new service
resource "layer0_service" "guestbook" {
  name          = "guestbook"
  environment   = "environment123"
  deploy        = "deploy123"
  load_balancer = "loadbalancer123"
  scale         = 3
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service
* `environment` - (Required) The id of the environment to place the service inside of
* `deploy` - (Required) The id of the deploy for the service to run
* `load_balancer` (Optional) The id of the load balancer to place the service behind
* `scale` (Optional, Default: 1) The number of copies of the service to run

### Attribute Reference
The following attributes are exported:

* `id` - The id of the service
* `name` - The name of the service
* `environment` - The id of the environment the service exists in
* `deploy` - The id of the deploy the service is running
* `load_balancer` - The id of the load balancer the service is behind (if `load_balancer` was set)
* `scale` - The current desired scale of the service

---

## Best Practices

* Always run `Terraform plan` before `terraform apply`.
This will show you what action(s) Terraform plans to make before actually executing them.
* Use [variables](https://www.terraform.io/intro/getting-started/variables.html) to reference secrets.
Secrets can be placed in a file named `terraform.tfvars`, or by setting `TF_VAR_*` environment variables.
More information can be found [here](https://www.terraform.io/intro/getting-started/variables.html).
* Use Terraform's `remote` command to backup and sync your `terraform.tfstate` file across different members in your organization.
Terraform has documentation for using S3 as a backend [here](https://www.terraform.io/docs/backends/types/s3.html).
* Terraform [modules](https://www.terraform.io/intro/getting-started/modules.html) allow you to define and consume reusable components.
* Example configurations can be found [here](https://github.com/hashicorp/Terraform/tree/master/examples)
