# Deployment 1: A Simple Guestbook App

In this section you'll learn how different Layer0 commands work together to deploy applications to the cloud.
The example application in this section is a guestbook -- a web application that acts as a simple message board.
You can choose to complete this section using either [the Layer0 CLI](#deploy-with-layer0-cli) or [Terraform](#deploy-with-terraform).


---

## Deploy with Layer0 CLI

If you're following along, you'll want to be working in the `walkthrough/deployment-1/` directory of your clone of the [guides](https://github.com/quintilesims/guides) repo.

Files used in this deployment:

| Filename | Purpose |
|----------|---------|
| `Guestbook.Dockerrun.aws.json` | Template for running the Guestbook application |


---

### Part 1: Create the Environment

The first step in deploying an application with Layer0 is to create an environment.
An environment is a dedicated space in which one or more services can reside.
Here, we'll create a new environment named **demo-env**.
At the command prompt, execute the following:

`l0 environment create demo-env`

We should see output like the following:

```
ENVIRONMENT ID  ENVIRONMENT NAME  CLUSTER COUNT  INSTANCE SIZE
demo00e6aa9     demo-env          0              m3.medium
```

We can inspect our environments in a couple of different ways:

- `l0 environment list` will give us a brief summary of all environments:

```
ENVIRONMENT ID  ENVIRONMENT NAME
demo00e6aa9     demo-env
api             api
```

- `l0 environment get demo-env` will show us more information about the **demo-env** environment we just created:

```
ENVIRONMENT ID  ENVIRONMENT NAME  CLUSTER COUNT  INSTANCE SIZE
demo00e6aa9     demo-env          0              m3.medium
```

- `l0 environment get \*` illustrates wildcard matching (you could also have used `demo*` in the above command), and will return detailed information for _each_ environment, not just one - it's like a detailed `list`:

```
ENVIRONMENT ID  ENVIRONMENT NAME  CLUSTER COUNT  INSTANCE SIZE
demo00e6aa9     demo-env          0              m3.medium
api             api               2              m3.medium
```

---

### Part 2: Create the Load Balancer

In order to expose a web application to the public internet, we need to create a load balancer.
A load balancer listens for web traffic at a specific address and directs that traffic to a Layer0 service.

A load balancer also has a notion of a health check - a way to assess whether or not the service is healthy and running properly.
By default, Layer0 configures the health check of a load balancer based upon a simple TCP ping to port 80 every thirty seconds.
Also by default, this ping will timeout after five seconds of no response from the service, and two consecutive successes or failures are required for the service to be considered healthy or unhealthy.

Here, we'll create a new load balancer named **guestbook-lb** inside of our environment named **demo-env**.
The load balancer will listen on port 80, and forward that traffic along to port 80 in the Docker container using the HTTP protocol.
Since the port configuration is already aligned with the default health check, we don't need to specify any health check configuration when we create this load balancer.
At the command prompt, execute the following:

`l0 loadbalancer create --port 80:80/http demo-env guestbook-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS       PUBLIC  URL
guestbodb65a     guestbook-lb       demo-env              80:80/HTTP  true
```

The following is a summary of the arguments passed in the above command:

* `loadbalancer create`: creates a new load balancer
* `--port 80:80/HTTP`: instructs the load balancer to forward requests from port 80 on the server to port 80 in the Docker container using the HTTP protocol
* `demo-env`: the name of the environment in which you are creating the load balancer
* `guestbook-lb`: a name for the load balancer itself

You can inspect load balancers in the same way that you inspected environments in Part 1.
Try running the following commands to get an idea of the information available to you:

- `l0 loadbalancer list`
- `l0 loadbalancer get guestbook-lb`
- `l0 loadbalancer get gues*`
- `l0 loadbalancer get \*`

!!! Note
	Notice that the load balancer `list` and `get` outputs list an `ENVIRONMENT` field - if you ever have load balancers (or other Layer0 entities) with the same name but in different environments, you can target a specific load balancer by qualifying it with its environment name:

	```
	`l0 loadbalancer get demo-env:guestbook-lb`
	```

---

### Part 3: Deploy the Docker Task Definition

The `deploy` command is used to specify the Docker task definition that refers to a web application.

Here, we'll create a new deploy called **guestbook-dpl** that refers to the **Guestbook.Dockerrun.aws.json** file found in the guides reposiory.
At the command prompt, execute the following:

`l0 deploy create Guestbook.Dockerrun.aws.json guestbook-dpl`

We should see output like the following:

```
DEPLOY ID        DEPLOY NAME    VERSION
guestbook-dpl.1  guestbook-dpl  1
```

The following is a summary of the arguments passed in the above command:

* `deploy create`: creates a new deployment and allows you to specify a Docker task definition
* `Guestbook.Dockerrun.aws.json`: the file name of the Docker task definition (use the full path of the file if it is not in your current working directory)
* `guestbook-dpl`: a name for the deploy, which you will use later when you create the service

!!! Note
	The `DEPLOY NAME` and `VERSION` are combined to create a unique identifier for a deploy. If you create additional deploys named **guestbook-dpl**, they will be assigned different version numbers.

	You can always specify the latest version when targeting a deploy by using `<deploy name>:latest` -- for example, `guestbook-dpl:latest`.

Deploys support the same methods of inspection as environments and load balancers:

- `l0 deploy list`
- `l0 deploy get guestbook*`
- `l0 deploy get \*`


---

### Part 4: Create the Service

The final stage of the deployment process involves using the `service` command to create a new service and associate it with the environment, load balancer, and deploy that we created in the previous sections.
The service will execute the Docker containers which have been described in the deploy.

Here, we'll create a new service called **guestbook-svc**. At the command prompt, execute the following:

`l0 service create --loadbalancer demo-env:guestbook-lb demo-env guestbook-svc guestbook-dpl:latest`

We should see output like the following:

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:1*  0/1
```

The following is a summary of the arguments passed in the above command:

* `service create`: creates a new service
* `--loadbalancer demo-env:guestbook-lb`: the fully-qualified name of the load balancer; in this case, the load balancer named **guestbook-lb** in the environment named **demo-env**. 
	- _(It is not strictly necessary to use the fully qualified name of the load balancer, unless another load balancer with exactly the same name exists in a different environment.)_
* `demo-env`: the name of the environment you created in Part 1
* `guestbook-svc`: a name for the service you are creating
* `guestbook-dpl`: the name of the deploy that you created in Part 3

Layer0 services can be queried using the same `get` and `list` commands that we've come to expect by now.


---

### Check the Status of the Service

After a service has been created, it may take several minutes for that service to completely finish deploying.
A service's status may be checked by using the `service get` command.

Let's take a peek at our **guestbook-svc** service.
At the command prompt, execute the following:

`l0 service get demo-env:guestbook-svc`

If we're quick enough, we'll be able to see the first stage of the process (this is what was output after running the `service create` command up in Part 4).
We should see an asterisk (\*) next to the name of the **guestbook-dpl:1** deploy, which indicates that the service is in a transitional state:
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:1*  0/1
```

In the next phase of deployment, if we execute the `service get` command again, we will see **(1)** in the **Scale** column; this indicates that 1 copy of the service is transitioning to an active state:
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:1*  0/1 (1)
```

We should see output like the following:

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:1   1/1
```

!!! Note
	More detailed information about the state of a service may be acquired by running the following command:

	`l0 service logs <SERVICE>`


---

### Get the Application's URL

Once the service has been completely deployed, we can obtain the URL for the application and launch it in a browser.

At the command prompt, execute the following:

`l0 loadbalancer get demo-env:guestbook-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE        PORTS       PUBLIC  URL
guestbodb65a     guestbook-lb       demo-env     guestbook-svc  80:80/HTTP  true    <url>
```

Copy the value shown in the **URL** column and paste it into a web browser.
The guestbook application will appear (once the service has completely finished deploying).


---

### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Layer0 to delete the environment and terminate the application.

`l0 environment delete demo-env`

However, if you intend to continue through [Deployment 2](deployment-2), you will want to keep the resources you made in this section.


---

## Deploy with Terraform

Instead of using the Layer0 CLI directly, you can instead use our Terraform provider, and deploy using Terraform _([learn more](/reference/terraform-plugin))_.
You can use Terraform with Layer0 and AWS to create "fire-and-forget" deployments for your applications.

If you're following along, you'll want to be working in the `walkthrough/deployment-1/` directory of your clone of the [guides](https://github.com/quintilesims/guides) repo.

We use these files to set up a Layer0 environment with Terraform:

| Filename | Purpose |
|----------|---------|
| `main.tf` | Provisions resources; populates resources in template files |
| `outputs.tf` | Values that Terraform will yield during deployment |
| `terraform.tfstate` | Tracks status of deployment _(created and managed by Terraform)_ |
| `terraform.tfvars` | Variables specific to the environment and application(s) |
| `variables.tf` | Variables that Terraform will use during deployment |

### `*.tf`: A Brief Aside

Let's take a moment to discuss the `.tf` files.
The names of these files (and even the fact that they are separated out into multiple files at all) are completely arbitrary and exist soley for human-readability.
Terraform understands all `.tf` files in a directory all together.

In `variables.tf`, you'll see `"endpoint"` and `"token"` variables.

In `outputs.tf`, you'll see that Terraform should spit out the url of the guestbook's load balancer once deployment has finished.

In `main.tf`, you'll see the bulk of the deployment process.
If you've followed along with the Layer0 CLI deployment above, it should be fairly easy to see how blocks in this file map to steps in the CLI process.
When we began the CLI deployment, our first step was to create an environment:

`l0 environment create demo-env`

This command is recreated in `main.tf` like so:

```hcl
# walkthrough/deployment-1/main.tf

resource "layer0_environment" "demo-env" {
	name = "demo-env"
}
```

We've bundled up the heart of the Guestbook deployment (load balancer, deploy, service, etc.) into a [Terraform module](https://www.terraform.io/docs/modules/usage.html).
To use it, we declare a `module` block and pass in the source of the module as well as any configuration or variables that the module needs.

```hcl
# walkthrough/deployment-1/main.tf

module "guestbook" {
    source         = "github.com/quintilesims/guides//guestbook/module"
    environment_id = "${layer0_environment.demo.id}"
}
```

You can see that we pass in the ID of the environment we create.
All variables declared in this block are passed to the module, so the next file we should look at is `variables.tf` inside of the module to get an idea of what the module is expecting.

There are a lot of variables here, but only one of them doesn't have a default value.

```hcl
# guestbook/module/variables.tf

variable "environment_id" {
    description = "id of the layer0 environment in which to create resources"
}
```

You'll notice that this is the variable that we're passing in.
For this particular deployment of the Guestbook, all of the default options are fine.
We could override any of them if we wanted to, just by specifying a new value for them back in `deployment-1/main.tf`.

Now that we've seen the variables that the module will have, let's take a look at part of `module/main.tf` and see how some of them might be used:

```hcl
# guestbook/module/main.tf

resource "layer0_load_balancer" "guestbook-lb" {
	name = "${var.load_balancer_name}"
	environment = "${var.environment_id}"
	port {
		host_port = 80
		container_port = 80
		protocol = "http"
	}
}

...
```

You can follow [this link](/reference/terraform-plugin/) to learn more about Layer0 resources in Terraform.


---

### Part 1: Terraform Get

This deployment uses modules, so we'll need to fetch those source materials.
At the command prompt, execute the following command:

`terraform get`

We should see output like the following:

```
Get: git::https://github.com/quintilesims/guides.git
```

We should now have a new local directory called `.terraform/`.
We don't need to do anything with it; we just want to make sure it's there.


---

### Part 2: Terraform Plan

Before we actually create/update/delete any resources, it's a good idea to find out what Terraform intends to do.

Run `terraform plan`. Terraform will prompt you for configuration values that it does not have:

```
var.endpoint
	Enter a value:

var.token
	Enter a value:
```
You can find these values by running `l0-setup endpoint <your layer0 prefix>`.

!!! Note
	There are a few ways to configure Terraform so that you don't have to keep entering these values every time you run a Terraform command (editing the `terraform.tfvars` file, or exporting evironment variables like `TF_VAR_endpoint` and `TF_VAR_token`, for example). See the [Terraform Docs](https://www.terraform.io/docs/configuration/variables.html) for more.

The `plan` command should give us output like the following:

```
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.template_file.guestbook: Refreshing state...
The Terraform execution plan has been generated and is shown below.
Resources are shown in alphabetical order for quick scanning. Green resources
will be created (or destroyed and then created if an existing resource
exists), yellow resources are being changed in-place, and red resources
will be destroyed. Cyan entries are data sources to be read.

Note: You didn't specify an "-out" parameter to save this plan, so when
"apply" is called, Terraform can't guarantee this is what will execute.

+ layer0_environment.demo
    ami:               "<computed>"
    cluster_count:     "<computed>"
    links:             "<computed>"
    name:              "demo"
    os:                "linux"
    security_group_id: "<computed>"
    size:              "m3.medium"

+ module.guestbook.layer0_deploy.guestbook
    content: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"guestbook\",\n            \"image\": \"quintilesims/guestbook\",\n            \"essential\": true,\n      \"memory\": 128,\n            \"environment\": [\n                {\n                    \"name\": \"GUESTBOOK_BACKEND_TYPE\",\n                    \"value\": \"memory\"\n                },\n                {\n          \"name\": \"GUESTBOOK_BACKEND_CONFIG\",\n                    \"value\": \"\"\n                },\n           {\n                    \"name\": \"AWS_ACCESS_KEY_ID\",\n                    \"value\": \"\"\n        },\n                {\n                    \"name\": \"AWS_SECRET_ACCESS_KEY\",\n                    \"value\": \"\"\n                },\n                {\n                    \"name\": \"AWS_REGION\",\n      \"value\": \"us-west-2\"\n                }\n            ],\n            \"portMappings\": [\n   {\n                    \"hostPort\": 80,\n                    \"containerPort\": 80\n                }\n      ]\n        }\n    ]\n}\n"
    name:    "guestbook"

+ module.guestbook.layer0_load_balancer.guestbook
    environment:                    "${var.environment_id}"
    health_check.#:                 "<computed>"
    name:                           "guestbook"
    port.#:                         "1"
    port.2027667003.certificate:    ""
    port.2027667003.container_port: "80"
    port.2027667003.host_port:      "80"
    port.2027667003.protocol:       "http"
    url:                            "<computed>"

+ module.guestbook.layer0_service.guestbook
    deploy:        "${ var.deploy_id == \"\" ? layer0_deploy.guestbook.id : var.deploy_id }"
    environment:   "${var.environment_id}"
    load_balancer: "${layer0_load_balancer.guestbook.id}"
    name:          "guestbook"
    scale:         "1"
    wait:          "true"


Plan: 4 to add, 0 to change, 0 to destroy.
```

This shows you that Terraform intends to create a deploy, an environment, a load balancer, and a service, all through Layer0.

If you've gone through this deployment using the [Layer0 CLI](#deploy-with-layer0-cli), you may notice that these resources appear out of order - that's fine. Terraform presents these resources in alphabetical order, but underneath, it knows the correct order in which to create them.

Once we're satisfied that Terraform will do what we want it to do, we can move on to actually making these things exist!


---

### Part 3: Terraform Apply

Run `terraform apply` to begin the process.

We should see output like the following:

```
layer0_environment.demo: Refreshing state...
...
...
...
layer0_service.guestbook: Creation complete

Apply complete! Resources: 7 added, 0 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: terraform.tfstate

Outputs:

guestbook_url = <http endpoint for the sample application>
```

!!! Note
	It may take a few minutes for the guestbook service to launch and the load balancer to become available. During that time you may get HTTP 503 errors when making HTTP requests against the load balancer URL.


### What's Happening

Terraform provisions the AWS resources through Layer0, configures environment variables for the application, and deploys the application into a Layer0 environment. Terraform also writes the state of your deployment to the `terraform.tfstate` file (creating a new one if it's not already there).


### Cleanup

When you're finished with the example, you can instruct Terraform to destroy the Layer0 environment, and terminate the application. Execute the following command (in the same directory):

`terraform destroy`

Then, remove the `.terraform/` directory.


---

