# An Iterative Walkthrough

This guide aims to take you through three increasingly-complex deployment examples using Layer0. Successive sections build upon the previous ones, and each deployment can be completed either through using the Layer0 CLI directly, or through Terraform using our custom [Layer0 Terraform Provider](/reference/terraform-plugin).

We assume that you're using Layer0 v0.9.0 or later. If you have not already installed and configured Layer0, see the [installation guide](/setup/install). If you are running an older version of Layer0, you may need to [upgrade](/setup/upgrade#upgrading-older-versions-of-layer0).

If you intend to deploy services using the Layer0 Terraform Provider, you'll want to make sure that you've [installed](/reference/terraform-plugin/#install) the provider correctly.

Regardless of the deployment method you choose, you should also clone/download our [examples repository](https://github.com/quintilesims/layer0-examples/), which contains all the files you will need to progress through our guides. As you do so, we will assume that your working directory matches the part of the guide that you're following (for example, Deployment 1 of this guide will assume that your working directory is `.../layer0-examples/iterative-walkthrough/deployment-1/`).

**Table of Contents**:

- [Deployment 1](#deployment-1-a-simple-guestbook-app): Deploying a web service (Guestbook)
- [Deployment 2](#deployment-2-guestbook-redis): Deploying Guestbook and a data store service (Redis)
- [Deployment 3](#deployment-3-guestbook-redis-consul): Deploying Guestbook, Redis, and a service discovery service (Consul)


## Deployment 1: A Simple Guestbook App

In this section you'll learn how different Layer0 commands work together to deploy applications to the cloud. The example application in this section is a guestbook -- a web application that acts as a simple message board. You can choose to complete this section using either the [Layer0 CLI](#1a-deploy-with-layer0-cli) or [Terraform](#1b-deploy-with-terraform).


## 1a: Deploy with Layer0 CLI

Remember, we assume that you're working in the `iterative-walkthrough/deployment-1/` directory, which contains the files you need for some of the commands in this section to work.


### Part 1: Create the Environment

The first step in deploying an application with Layer0 is to create an environment. An environment is a dedicated space in which one or more services can reside. Here, we'll create a new environment named **demo-env**. At the command prompt, execute the following:

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

In order to expose a web application to the public internet, we need to create a load balancer. A load balancer listens for web traffic at a specific address and directs that traffic to a Layer0 service.

A load balancer also has a notion of a health check - a way to assess whether or not the service is healthy and running properly. By default, Layer0 configures the health check of a load balancer based upon a simple TCP ping to port 80 every thirty seconds. Also by default, this ping will timeout after five seconds of no response from the service, and two consecutive successes or failures are required for the service to be considered healthy or unhealthy.

Here, we'll create a new load balancer named **guestbook-lb** inside of our environment named **demo-env**. The load balancer will listen on port 80, and forward that traffic along to port 80 in the Docker container using the HTTP protocol. Since the port configuration is already aligned with the default health check, we don't need to specify any health check configuration when we create this load balancer. At the command prompt, execute the following:

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

You can inspect load balancers in the same way that you inspected environments in Part 1. Try running the following commands to get an idea of the information available to you:

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

Here, we'll create a new deploy called **guestbook-dpl** that refers to the **Guestbook.Dockerrun.aws.json** file you obtained earlier. At the command prompt, execute the following:

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

The final stage of the deployment process involves using the `service` command to create a new service and associate it with the environment, load balancer, and deploy that we created in the previous sections. The service will execute the Docker containers which have been described in the deploy.

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

After a service has been created, it may take several minutes for that service to completely finish deploying. A service's status may be checked by using the `service get` command.

Let's take a peek at our **guestbook-svc** service. At the command prompt, execute the following:

`l0 service get demo-env:guestbook-svc`

If we're quick enough, we'll be able to see the first stage of the process (this is what was output after running the `service create` command up in Part 4). We should see an asterisk (\*) next to the name of the **guestbook-dpl:1** deploy, which indicates that the service is in a transitional state:
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

Copy the value shown in the **URL** column and paste it into a web browser. The guestbook application will appear (once the service has completely finished deploying).


---

### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Layer0 to delete the environment and terminate the application.

`l0 environment delete demo-env`

However, if you intend to continue through [Deployment 2](#deployment-2-guestbook-redis), you will want to keep the resources you made in this section.


---

## 1b: Deploy with Terraform

Instead of using the Layer0 CLI directly, you can instead use our Terraform provider, and deploy using Terraform _([learn more](/reference/terraform-plugin))_. You can use Terraform with Layer0 and AWS to create "fire-and-forget" deployments for your applications.

Remember, we assume that you've cloned the [layer0-examples](https://github.com/quintilesims/layer0-examples) repo and are working in the `iterative-walkthrough/deployment-1/` directory.

We use these files to set up a Layer0 environment with Terraform:

|Filename|Purpose|
|---|---|
|`Guestbook.Dockerrun.aws.json`|Template for running the Guestbook application|
|`layer0.tf`|Provisions resources; populates resources in template files|
|`terraform.tfstate`|Tracks status of deployment _(created and managed by Terraform)_|
|`terraform.tfvars`|Variables specific to the environment and application(s)|

### `Layer0.tf`: A Brief Aside

Let's take a moment to look through the `layer0.tf` file. If you've followed along with the Layer0 CLI deployment above, it should be fairly easy to see how blocks in this file map to steps in the CLI process.

When we began the CLI deployment, our first step was to create an environment:

`l0 environment create demo-env`

This command is recreated in `layer0.tf` like so:

```
resource "layer0_environment" "demo-env" {
	name = "demo-env"
}
```

The value of the `name` field inside the resource block is the name that maps to `demo-env` in the Layer0 CLI command. The `"demo-env"` in the resource declaration line is the identifier that Terraform will use to reference this resource.

The next step in the CLI process would be to make a load balancer:

`l0 loadbalancer create --port 80:80/http demo-env guestbook-lb`

In `layer0.tf`:

```
resource "layer0_load_balancer" "guestbook-lb" {
	name = "guestbook-lb"
	environment = "${layer0_environment.demo-env.id}"
	port {
		host_port = 80
		container_port = 80
		protocol = "http"
	}
}
```

We use Terraform's interpolation syntax to discover and use the ID of our environment. The format is pretty simple:

`${<resource_type>.<resource_identifier>.<property>}`

And that's about all you need to be able to understand the file! You can follow [this link](/reference/terraform-plugin/) to learn more about Layer0 resources in Terraform.

---

### Part 1: Terraform Plan

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
+ layer0_deploy.guestbook
    content: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"guestbook\",\n            \"image\": \"quintilesims/guestbook\",\n            \"essential\": true,\n            \"memory\": 128,\n            \"portMappings\": [\n                {\n                    \"hostPort\": 80,\n                    \"containerPort\": 80\n                }\n            ]\n        }\n    ]\n}\n"
    name:    "guestbook"

+ layer0_environment.demo
    cluster_count:     "<computed>"
    name:              "demo"
    security_group_id: "<computed>"
    size:              "m3.medium"

+ layer0_load_balancer.guestbook
    environment:                    "${layer0_environment.demo.id}"
    health_check.#:                 "<computed>"
    name:                           "guestbook"
    port.#:                         "1"
    port.2027667003.certificate:    ""
    port.2027667003.container_port: "80"
    port.2027667003.host_port:      "80"
    port.2027667003.protocol:       "http"
    url:                            "<computed>"

+ layer0_service.guestbook
    deploy:        "${layer0_deploy.guestbook.id}"
    environment:   "${layer0_environment.demo.id}"
    load_balancer: "${layer0_load_balancer.guestbook.id}"
    name:          "guestbook"
    scale:         "1"


Plan: 4 to add, 0 to change, 0 to destroy.
```

This shows you that Terraform intends to create a deploy, an environment, a load balancer, and a service, all through Layer0.

If you've gone through [Deployment 1a](#1a-deploy-with-layer0-cli) which used the Layer0 CLI, you may notice that these resources appear out of order - that's fine. Terraform presents these resources in alphabetical order, but underneath, it knows the correct order in which to create them.

Once we're satisfied that Terraform will do what we want it to do, we can move on to actually making these things exist!


---

### Part 2: Terraform Apply

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


---

## Deployment 2: Guestbook + Redis

In this section, we're going to add some complexity to the previous deployment. [Deployment 1](#deployment-1-a-simple-guestbook-app) saw us create a simple guestbook application which kept its data in memory. But what if that ever came down, either by intention or accident? It would be easy enough to redeploy it, but all of the entered data would be lost. For this deployment, we're going to separate the data store from the guestbook application by creating a second Layer0 service which will house a local Redis database server and linking it to the first. You can choose to complete this section using either the [Layer0 CLI](#2a-deploy-with-layer0-cli) or [Terraform](#2b-deploy-with-terraform).


## 2a: Deploy with Layer0 CLI

For this example, we'll be working in the `iterative-walkthrough/deployment-2/` directory.


### Part 1: Create the Redis Load Balancer

Both the Guestbook service and the Redis service will live in the same Layer0 environment, so we don't need to create one like we did in the first deployment. We'll start by making a load balancer behind which the Redis service will be deployed.

The `Redis.Dockerrun.aws.json` task definition file we'll use is very simple - it just spins up a Redis server with the default configuration, which means that it will be serving on port 6379. Our load balancer needs to be able to forward TCP traffic to and from this port. And since we don't want the Redis server to be exposed to the public internet, we'll put it behind a private load balancer; private load balancers only accept traffic that originates from within their own environment. At the command prompt, execute the following:

`l0 loadbalancer create --port 6379:6379/tcp --private demo-env redis-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS          PUBLIC  URL
redislb16ae6     redis-lb           demo-env              6378:6379:TCP  false
```

The following is a summary of the arguments passed in the above command:

- `loadbalancer create`: creates a new load balancer
- `--port 6379:6379/TCP`: instructs the load balancer to forward requests from port 6379 on the load balancer to port 6379 in the EC2 instance using the TCP protocol
- `--private`: instructs the load balancer to ignore external traffic
- `demo-env`: the name of the environment in which the load balancer is being created
- `redis-lb`: a name for the load balancer itself


---

### Part 2: Deploy the Docker Task Definition

Here, we just need to create the deploy using the `Redis.Dockerrun.aws.json` task definition file. At the command prompt, execute the following:

`l0 deploy create Redis.Dockerrun.aws.json redis-dpl`

We should see output like the following:

```
DEPLOY ID    DEPLOY NAME  VERSION
redis-dpl.1  redis-dpl    1
```

The following is a summary of the arguments passed in the above command:

- `deploy create`: creates a new Layer0 Deploy and allows you to specify a Docker task definition
- `Redis.Dockerrun.aws.json`: the file name of the Docker task definition (use the full path of the file if it is not in your current working directory)
- `redis-dpl`: a name for the deploy, which we will use later when we create the service


---

### Part 3: Create the Redis Service

Here, we just need to pull the previous resources together into a service. At the command prompt, execute the following:

`l0 service create --loadbalancer demo-env:redis-lb demo-env redis-svc redis-dpl:latest`

We should see output like the following:

```
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
redislb16ae6  redis-svc     demo-env     redis-lb      redis-dpl:1  0/1
```

The following is a summary of the arguments passed in the above commands:

- `service create`: creates a new Layer0 Service
- `--loadbalancer demo-env:redis-lb`: the fully-qualified name of the load balancer; in this case, the load balancer named **redis-lb** in the environment named **demo-env**
    - _(Again, it's not strictly necessary to use the fully-qualified name of the load balancer as long as there isn't another load balancer with the same name)_
- `demo-env`: the name of the environment in which the service is to reside
- `redis-svc`: a name for the service we're creating
- `redis-dpl:latest`: the name of the deploy the service will put into action
    - _(We use `:` to specify which deploy we want - `:latest` will always give us the most recently-created one.)_


---

### Part 4: Check the Status of the Redis Service

As in the first deployment, we can keep an eye on our service by using the `service get` command:

`l0 service get redis-svc`

Once the service has finished scaling, try looking at the service's logs to see the output that the Redis server creates:

`l0 service logs redis-svc`

Among some warnings and information not important to this exercise and a fun bit of ASCII art, you should see something like the following:

```
... # words and ASCII art
1:M 05 Apr 23:29:47.333 * The server is now ready to accept connections on port 6379
```

Now we just need to teach the Guestbook application how to talk with our Redis service.


---

### Part 5: Update the Guestbook Deploy

You should see in `iterative-walkthrough/deployment-2/` another `Guestbook.Dockerrun.aws.json` file. This file is very similar to but not the same as the one in `deployment-1/` - if you open it up, you can see the following additions:

```
    ...
    "environment": [
        {
            "name": "REDIS_ADDRESS_AND_PORT",
            "value": "${redis_address}"
        }
    ],
    ...
```

That `value` is what will point the Guestbook application towards the Redis server. The `${redis_address}` needs to be replaced and populated in the following format:

```
"value": "ADDRESS_TO_REDIS_SERVER:PORT_THE_SERVER_IS_SERVING_ON"
```

We already know that Redis is serving on port 6379, so let's go find the server's address. Remember, it lives behind a load balancer that we made, so run the following command:

`l0 loadbalancer get redis-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE    PORTS          PUBLIC  URL
redislb16ae6     redis-lb           demo-env     redis-svc  6379:6379/TCP  false   internal-l0-<yadda-yadda>.elb.amazonaws.com
```

Copy that `URL` value, replace `${redis_address}` with the `URL` value in `Guestbook.Dockerrun.aws.json`, append `:6379` to it, and save the file. It should look something like the following:

```
    ...
    "environment": [
        {
            "name": "REDIS_ADDRESS_AND_PORT",
            "value": "internal-l0-<yadda-yadda>.elb.amazonaws.com:6379"
        }
    ],
    ...
```

Now, we can create an updated deploy:

`l0 deploy create Guestbook.Dockerrun.aws.json guestbook-dpl`

We should see output like the following:

```
DEPLOY ID        DEPLOY NAME    VERSION
guestbook-dpl.2  guestbook-dpl  2
```


---

### Part 6: Update the Guestbook Service

Almost all the pieces are in place! Now we just need to apply the new Guestbook deploy to the running Guestbook service:

`l0 service update guestbook-svc guestbook-dpl:latest`

As the Guestbook service moves through the phases of its update process, we should see outputs like the following (if we keep an eye on the service with `l0 service get guestbook-svc`, that is):

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2*  1/1
                                                        guestbook-dpl:1
```

_above: `guestbook-dpl:2` is in a transitional state_

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS      SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2  2/1
                                                        guestbook-dpl:1
```

_above: both versions of the deployment are running at scale_

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2   1/1
                                                        guestbook-dpl:1*
```
_above: `guestbook-dpl:1` is in a transitional state_
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS      SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2  1/1
```
_above: `guestbook-dpl:1` has been removed, and only `guestbook-dpl:2` remains_


---

### Part 7: Prove It

You should now be able to point your browser at the URL for the Guestbook loadbalancer (run `l0 loadbalancer get guestbook-svc` to find it) and see what looks like the same Guestbook application you deployed in the first section of the walkthrough. Go ahead and add a few entries, make sure it's functioning properly. We'll wait.

Now, let's prove that we've actually separated the data from the application by deleting and redeploying the Guestbook application:

`l0 service delete --wait guestbook-svc`
`l0 loadbalancer delete --wait guestbook-lb`

_(We'll leave the `deploy` intact so we can spin up a new service easily, and we'll leave the environment untouched because it also contained the Redis server. We'll also pass the `--wait` flag so that we don't need to keep checking on the status of the job to know when it's complete.)_

Once those resources have been deleted, we can recreate them!

Create another load balancer:

`l0 loadbalancer create --ports 80:80/http demo-env guestbook-lb`

Create another service, using the **guestbook-dpl** deploy we kept around:

`l0 service create --loadbalancer demo-env:guestbook-lb demo-env guestbook-svc guestbook-dpl:latest`

Wait for everything to spin up, and hit that new load balancer's url (`l0 loadbalancer get guestbook-lb`) with your browser. Your data should still be there!


---

### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Layer0 to delete the environment and terminate the application.

`l0 environment delete demo-env`

However, if you intend to continue through [Deployment 3](#deployment-3-guestbook-redis-consul), you will want to keep the resources you made in this section.


---

## 2b: Deploy with Terraform

As before, we can complete this deployment using Terraform and the Layer0 provider instead of the Layer0 CLI. As before, we will assume that you've cloned the [layer0-examples](https://github.com/quintilesims/layer0-examples) repo and are working in the `iterative-walkthrough/deployment-2/` directory.

We'll use these files to manage our deployment with Terraform:

|Filename|Purpose|
|---|---|
|`Guestbook.Dockerrun.aws.json`|Template for running the Guestbook application|
|`layer0.tf`|Provisions resources; populates variables in template files|
|`Redis.Dockerrun.aws.json`|Template for running the Redis application|
|`terraform.tfstate`|Tracks status of deployment _(created and managed by Terraform)_|
|`terraform.tfvars`|Variables specific to the environment and application(s)|


---

### `Layer0.tf`: A Brief Aside: Revisited

This file hasn't changed much since the first deployment - although, we _have_ added some things to it. You can see that the configurations for both the Guestbook _and_ the Redis services are contained within this single file.

If you've gone through the [2a Deployment](#2a-deploy-with-layer0-cli) using the Layer0 CLI above, you'll recall that we had to obtain the URL of the Redis load balancer and manually plug it into the `Guestbook.Dockerrun.aws.json` task definition file. You may be wondering how we're supposed to do that, when Terraform creates all the resources at once. Good news! We don't have to! We can configure `layer0.tf` to get that information and plug it in automatically!

When we manually edited that Guestbook task definition, we reconfigured a section that looked like this:

```
    ...
    "environment": [
        {
            "name": "REDIS_ADDRESS_AND_PORT",
            "value": "${redis_address}"
        }
    ],
    ...
```

You've probably noticed that `${redis_address}` is a Terraform [interpolation](https://www.terraform.io/docs/configuration/interpolation.html). If you look for the `template_file` resource labeled `guestbook` in `layer0.tf`, you'll see that in addition to supplying the Guestbook task definition file, we're also configuring a variable to be passed to the template:

```
data "template_file" "guestbook" {
	template = "${file("Guestbook.Dockerrun.aws.json")}"
    
    vars {
        redis_address = "${layer0_load_balancer.redis-lb.url}:6379"
    }
}
```

We find the URL of the Redis load balancer, append ":6379" to it, and use that to populate the `redis_address` variable in the task definition.


---

### Part 1: Terraform Plan

It's always a good idea to find out what Terraform intends to do, so let's do that:

`terraform plan`

As before, we'll be prompted for any variables Terraform needs and doesn't have (see the note in [Deployment 1b](#1b-deploy-with-terraform) for configuring Terraform variables). We'll see output similar to the following:

```
<= data.template_file.guestbook
    rendered: "<computed>"
    template: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"guestbook\",\n            \"image\": \"quintilesims/guestbook-redis\",\n            \"essential\": true,\n           \"memory\": 128,\n            \"environment\": [\n                {\n                    \"name\": \"REDIS_ADDRESS_AND_PORT\",\n                    \"value\": \"${redis_address}\"\n                }\n            ],\n           \"portMappings\": [\n                {\n                    \"hostPort\": 80,\n                    \"containerPort\": 80\n               }\n            ]\n        }\n    ]\n}\n"
    vars.%:   "<computed>"

+ layer0_deploy.guestbook-dpl
    content: "${data.template_file.guestbook.rendered}"
    name:    "guestbook-dpl"

+ layer0_deploy.redis-dpl
    content: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"redis\",\n            \"image\": \"redis:3.2-alpine\",\n            \"essential\": true,\n            \"memory\": 128,\n            \"portMappings\": [\n                {\n                    \"hostPort\": 6379,\n           \"containerPort\": 6379\n                }\n            ]\n        }\n    ]\n}\n"
    name:    "redis-dpl"

+ layer0_environment.demo-env
    cluster_count:     "<computed>"
    name:              "demo-env"
    security_group_id: "<computed>"
    size:              "m3.medium"

+ layer0_load_balancer.guestbook-lb
    environment:                    "${layer0_environment.demo-env.id}"
    health_check.#:                 "<computed>"
    name:                           "guestbook-lb"
    port.#:                         "1"
    port.2027667003.certificate:    ""
    port.2027667003.container_port: "80"
    port.2027667003.host_port:      "80"
    port.2027667003.protocol:       "http"
    url:                            "<computed>"

+ layer0_load_balancer.redis-lb
    environment:                    "${layer0_environment.demo-env.id}"
    health_check.#:                 "<computed>"
    name:                           "redis-lb"
    port.#:                         "1"
    port.1072619732.certificate:    ""
    port.1072619732.container_port: "6379"
    port.1072619732.host_port:      "6379"
    port.1072619732.protocol:       "tcp"
    private:                        "true"
    url:                            "<computed>"

+ layer0_service.guestbook-svc
    deploy:        "${layer0_deploy.guestbook-dpl.id}"
    environment:   "${layer0_environment.demo-env.id}"
    load_balancer: "${layer0_load_balancer.guestbook-lb.id}"
    name:          "guestbook-svc"
    scale:         "1"

+ layer0_service.redis-svc
    deploy:        "${layer0_deploy.redis-dpl.id}"
    environment:   "${layer0_environment.demo-env.id}"
    load_balancer: "${layer0_load_balancer.redis-lb.id}"
    name:          "redis-svc"
    scale:         "1"


Plan: 7 to add, 0 to change, 0 to destroy.
```

We should see that Terraform intends to add 7 new resources, some of which are for the Guestbook deployment and some of which are for the Redis deployment.


---

### Part 2: Terraform Apply

Run `terraform apply`, and we should see output similar to the following:

```
data.template_file.redis: Refreshing state...
layer0_deploy.redis-dpl: Creating...

...
...
...

layer0_service.guestbook-svc: Creation complete

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

!!! Note
	Again, Terraform writes the latest status of your deployment to `terraform.tfstate`. When you're finished with this section and ready to move on to [Deployment 3](#3b-deploy-with-terraform), destroy your Terraform deployment with `terraform destroy`.
    

---

## Deployment 3: Guestbook + Redis + Consul

In [Deployment 2](#2a-deploy-with-layer0-cli), we created two services in the same environment and linked them together manually.
While that can work for a small system, it's not really feasible for a system with a lot of moving parts - we would need to look up load balancer endpoints for all of our services and manually link them all together.
To that end, here we're going to to redeploy our two-service system using [Consul](https://www.consul.io), a service discovery tool.

For this deployment, we'll create a cluster of Consul servers which will be responsible for keeping track of the state of our system.
We'll also deploy new versions of the Guestbook and Redis task definition files - in addition to creating a container for its respective application, each task definition creates two other containers:

 - a container for a Consul agent, which is in charge of communicating with the Consul server cluster
 - a container for [Registrator](https://github.com/gliderlabs/registrator), which is charge of talking to the local Consul agent when a service comes up or goes down.

You can choose to complete this section using either the [Layer0 CLI](#3a-deploy-with-layer0-cli) or [Terraform](#3b-deploy-with-terraform).


## 3a: Deploy with Layer0 CLI

For this example, we'll be working in the `iterative-walkthrough/deployment-3/` directory.


---

### Part 1: Create the Consul Load Balancer

The Consul server cluster will live in the same environment as our Guestbook and Redis services - if you've completed the previous deployment, this environment already exists as **demo-env**.
We'll start by creating the load balancer behind which the Consul cluster will be deployed.
The load balancer is a private one, and is really only used to bootstrap the Consul servers into working order.
At the command prompt, execute the following:

`l0 loadbalancer create --port 8500:8500/tcp --port 8301:8301/tcp --private --healthcheck-target tcp:8500 demo-env consul-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS          PUBLIC  URL
consull66b23     consul-lb          consul-env            8500:8500/TCP  false
                                                          8301:8301/TCP
```

The following is a summary of the arguments passed in the above command:

- `loadbalancer create`: creates a new load balancer
- `--port 8500:8500/tcp`: instructs the load balancer to forward requests from port 8500 on the load balancer to port 8500 in the EC2 instance using the TCP protocol
- `--port 8301:8301/tcp`: instructs the load balancer to forward requests from port 8301 on the load balancer to port 8301 in the EC2 instance using the TCP protocol
- `--private`: instructs the load balancer to ignore outside traffic
- `--healthcheck-target`: instructs the load balancer to use a TCP ping on port 8500 as the basis for deciding whether the service is healthy
- `demo-env`: the name of the environment in which the load balancer is being created
- `consul-lb`: a name for the load balancer itself

While we're touching on the Consul load balancer, we should grab its URL - this is the one value that we'll need to know in order to deploy the rest of our system, no matter how large it may get.
At the command prompt, execute the following:

`l0 loadbalancer get consul-lb`

We should see output that looks like the output we just received above after creating the load balancer, but this time there is something in the **URL** column.
That URL is the value we're looking for.
Make note of it for when we reference it later.


---

### Part 2: Deploy the Consul Task Definition

Before we can create the deploy, we need to supply the URL of the Consul load balancer that we got in Part 1.
In `Consul.Dockerrun.aws.json`, find the entry in the `environment` block that looks like this:

```
{
    "name": "CONSUL_SERVER_URL",
    "value": "${consul_server_url}"
}
```

Replace `${consul_server_url}` with the Consul load balancer's URL and save the file.
We can then create the deploy.
At the command prompt, execute the following:

`l0 deploy create Consul.Dockerrun.aws.json consul-dpl`

We should see output like the following:

```
DEPLOY ID     DEPLOY NAME  VERSION
consul-dpl.1  consul-dpl   1
```

The following is a summary of the arguments passed in the above command:

- `deploy create`: creates a new Layer0 Deploy and allows you to specifiy a Docker task definition
- `Consul.Dockerrun.aws.json`: the file name of the Docker task definition (use the full path of the file if it is not in the current working directory)
- `consul-dpl`: a name for the deploy, which will later be used in creating the service


---

### Part 3: Create the Consul Service

Here, we pull the previous resources together to create a service.
At the command prompt, execute the following:

`l0 service create --wait --loadbalancer demo-env:consul-lb demo-env consul-svc consul-dpl:latest`

We should see output like the following:

```
Waiting for Deployment...
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS   SCALE
consuls2f3c6  consul-svc    demo-env     consul-lb     consul-dpl:1  1/1
```

The following is a summary of the arguments passed in the above commands:

- `service create`: creates a new Layer0 Service
- `--wait`: instructs the CLI to keep hold of the shell until the service has been successfully deployed
- `--loadbalancer demo-env:consul-lb`: the fully-qualified name of the load balancer behind which the service should live; in this case, the load balancer named **consul-lb** in the environment named **demo-env**
- `demo-env`: the name of the environment in which the service is to reside
- `consul-svc`: a name for the service itself
- `consul-dpl:latest`: the name and version of the deploy that the service should put into action

Once the service has finished being deployed (and `--wait` has returned our shell to us), we need to scale the service.

Currently, we only have one Consul server running in the cluster.
For best use, we should have at least 3 servers running (see [this link](https://www.consul.io/docs/internals/consensus.html) for more details on Consul servers and their concensus protocol).
Indeed, if we inspect the `command` block of the task definition file, we can find the following parameter: `-bootstrap-expect=3`.
This tells the Consul server that we have just deployed that it should be expecting a total of three servers.
We still need to fulfill that expectation, so we'll scale our service up to three.
At the command prompt, execute the following:

`l0 service scale --wait consul-svc 3`

We should see output like the following:

```
Waiting for Deployment...
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS   SCALE
consuls2f3c6  consul-svc    demo-env     consul-lb     consul-dpl:1  3/3
```

!!! WARNING "Important!"
    The successful completion of the **scale** command doesn't mean that we're ready to move on just yet!
    We need to check in on the logs (**l0 service logs consul-svc**) until we can confirm that all three of the Consul servers have synced up with each other.
    Each **consul-server** section in the logs should be ending with **consul: Adding LAN server [ip address]** or **agent: Join completed**.
    If you see one of the sections ending with **agent: Join failed, retrying in 30s**, you need to wait for that server to join the cluster before continuing.


---

### Part 4: Update and Redeploy the Redis and Guestbook Applications

We're going to need the URL of the Consul load balancer again.
In each of the Redis and Guestbook task definition files, look for the `CONSUL_SERVER_URL` block in the `consul-agent` container and replace the value field with the Consul load balancer URL, then save the file.
At the command prompt, execute the two following commands to create new version of the deploys for the Redis and Guestbook applications:

`l0 deploy create Redis.Dockerrun.aws.json redis-dpl`

`l0 deploy create Guestbook.Dockerrun.aws.json guestbook-dpl`

Then, execute the two following commands to redeploy the existing Redis and Guestbook services using those new deploys:

`l0 service update --wait redis-svc redis-dpl:latest`

`l0 service update --wait guestbook-svc guestbook-dpl:latest`

!!! NOTE
	Here, we should run `l0 service logs consul-svc` again and confirm that the Consul cluster has discovered these two services.

We can use `l0 loadbalancer get guestbook-lb` to obtain the guestbook application's URL, and then hit it with a web browser.
Our guestbook app should be up and running - this time, it's been deployed without needing to know the address of the Redis backend!

Of course, this is a simple example; in both this deployment and [Deployment 2](#2a-deploy-with-layer0-cli), we needed to use `l0 loadbalancer get` to obtain the URL of a load balancer.
However, in a system with many services that uses Consul like this example, we only ever need to find the URL of the Consul cluster - not the URLs of every service that needs to talk to another of our services.


---

### Part 5: Inspect the Consul Universe (Optional)

Let's take a glimpse into how this system that we've deployed works.
**This requires that we have access to the key pair we've told Layer0 about when we [set it up](/setup/install/#part-2-create-an-access-key).**


#### Open Ports for SSH

We want to SSH into the Guestbook EC2 instance, which means that we need to tell the Guestbook load balancer to allow SSH traffic through.
At the command prompt, execute the following:

`l0 loadbalancer addport guestbook-lb 22:22/tcp`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE        PORTS       PUBLIC  URL
guestbodb65a     guestbook-lb       demo-env     guestbook-svc  80:80/HTTP  true    <url>
                                                                22:22/TCP
```

We need to take note of the load balancer's URL here, too.


#### SSH Into the Instance

At the command prompt, execute the following:

`ssh -i /path/to/keypair ec2-user@<guestbook_load_balancer_url> -o ServerAliveInterval=30`

(We'll probably be asked if we want to continue connecting - we do, so we'll enter `yes`.)

Summary of arguments passed into the above command:

- `-i /path/to/keypair`: this allows us to specify an identity file for use when connecting to the remote machine - in this case, we want to replace `/path/to/keypair` with the actual path to the keypair we created when we set up Layer0
- `ec2-user@<guestbook_load_balancer_url>`: the address (here we want to replace `<guestbook_load_balancer_url>` with the actual URL of the guestbook load balancer) of the machine to which we want to connect and the name of the user (`ec2-user`) that we'd like to connect as
- `-o`: allows us to set parameters on the `ssh` command
- `ServerAliveInterval=30`: one of those `ssh` parameters - AWS imposes an automatic disconnect if a connection is not active for a certain amount of time, so we use this option to ping every 30 seconds to prevent that automatic disconnect


#### Look Around You

We're now inside of the EC2 instance!
If we run `docker ps`, we should see that our three Docker containers (the Guestbook app, a Consul agent, and Registrator) are up and running, as well as an `amazon-ecs-agent` image.
But that's not the Consul universe that we came here to see.
At the EC2 instance's command prompt, execute the following:

`echo $(curl -s localhost:8500/v1/catalog/services) | jq '.'`

We should see output like the following:

```
{
  "consul": [],
  "consul-8301": [
    "udp"
  ],
  "consul-8500": [],
  "consul-8600": [
    "udp"
  ],
  "guestbook-redis": [],
  "redis": []
}
```

Summary of commands passed in the above command:

- `curl -s http://localhost:8500/v1/catalog/services`: use `curl` to send a GET request to the specified URL, where `localhost:8500` is an HTTP connection to the local Consul agent in this EC2 instance (the `-s` flag just silences excess output from `curl`)
- `| jq '.'`: use a pipe (`|`) to take whatever returns from the left side of the pipe and pass it to the `jq` program, which we use here simply to pretty-print the JSON response
- `echo $(...)`: print out whatever returns from running the stuff inside of the parens; not necessary, but it gives us a nice newline after we get our response

In that output, we can see all of the things that our local Consul agent knows about.
In addition to a few connections to the Consul server cluster, we can see that it knows about the Guestbook application running in this EC2 instance, as well as the Redis application running in a different instance with its own Consul agent and Registrator.

Let's take a closer look at the Redis service and see how our Guestbook application is locating our Redis application.
At the EC2 instance's command prompt, execute the following:

`echo $(curl -s http://localhost:8500/v1/catalog/service/redis) | jq '.'`

We should see output like the following:

```
[
  {
    "ID": "b4bb81e6-fe6a-c630-2553-7f6492ae5275",
    "Node": "ip-10-100-230-97.us-west-2.compute.internal",
    "Address": "10.100.230.97",
    "Datacenter": "dc1",
    "TaggedAddresses": {
      "lan": "10.100.230.97",
      "wan": "10.100.230.97"
    },
    "NodeMeta": {},
    "ServiceID": "562aceee6935:ecs-l0-tlakedev-redis-dpl-20-redis-e0f989e5af97cdfd0e00:6379",
    "ServiceName": "redis",
    "ServiceTags": [],
    "ServiceAddress": "10.100.230.97",
    "ServicePort": 6379,
    "ServiceEnableTagOverride": false,
    "CreateIndex": 761,
    "ModifyIndex": 761
  }
]

```

Our Guestbook application makes a call like this one and figures out how to connect to the Redis service by mushing together the information from the `ServiceAddress` and `ServicePort` fields!

To close the `ssh` connection to the EC2 instance, run `exit` in the command prompt.


---

### Cleanup

When you're finished with the example, we can instruct Layer0 to terminate the applications and delete the environment.

`l0 environment delete demo-env`


---

## 3b: Deploy with Terraform

For this example, we'll be working in the `iterative-walkthrough/deployment-3/` directory.


---

### Part 1:


---

