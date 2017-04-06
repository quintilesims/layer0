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

We use these files to set up a Layer0 envrionment with Terraform:

|Filename|Purpose|
|----|----|
|`terraform.tfvars`|Variables specific to the environment and guestbook application|
|`Guestbook.Dockerrun.aws.json`|Template for running the guestbook application in a Layer0 environment|
|`layer0.tf`|Provision Layer0 resources and populate variables in `Guestbook.Dockerrun.aws.json`|


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


### What's happening

Terraform provisions the AWS resources through Layer0, configures environment variables for the application, and deploys the application into a Layer0 environment. Terraform also writes the state of your deployment to the `terraform.tfstate` file (creating a new one if it's not already there).


### Cleanup

When you're finished with the example, you can instruct Terraform to destroy the Layer0 environment, and terminate the application. Execute the following command (in the same directory):

`terraform destroy`

!!! Note
	As previously stated, Terraform writes the latest status of your deployment to `terraform.tfstate`. As you move on to [Deployment 2](#2b-deploy-with-terraform) of this walkthrough, it will be easier for you to just destroy your Terraform deployment with `terraform destroy`.

	However, the `apply` command is idempotent -- if you'd like, you may bring your `terraform.tfstate` file with you when you change directories to `iterative-walkthrough/deployment-2/`. When you run `terraform apply` in the next section of this walkthrough, Terraform will create new resources as expected - but it will also update existing resources that have changes and destroy resources that are no longer necessary, and update `terraform.tfstate` accordingly.


---

## Deployment 2: Guestbook + Redis

In this section, we're going to add some complexity to the previous deployment. [Deployment 1](#deployment-1-a-simple-guestbook-app) saw us create a simple guestbook application which kept its data in memory. But what if that ever came down, either by intention or accident? It would be easy enough to redeploy it, but all of the entered data would be lost. For this deployment, we're going to separate the data store from the guestbook application by creating a second Layer0 service which will house a local Redis database server and linking it to the first. You can choose to complete this section using either the [Layer0 CLI](#2a-deploy-with-layer0-cli) or [Terraform](#2b-deploy-with-terraform).


## 2a: Deploy with Layer0 CLI

For this example, we'll be working in the `iterative-walkthrough/deployment-2/` directory.


### Part 1: Create the Redis Load Balancer

Both the Guestbook service and the Redis service will live in the same Layer0 environment, so we don't need to create one like we did in the first deployment. We'll start by making a load balancer behind which the Redis service will be deployed.

The `Redis.Dockerrun.aws.json` task definition file we'll use is very simple - it just spins up a Redis server with the default configuration, which means that it will be serving on port 6379. Our load balancer needs to be able to forward TCP traffic to and from this port. And since the Redis server has no reason for anything other than the Guestbook application to be accessing it, we'll make this load balancer private. At the command prompt, execute the following:

`l0 loadbalancer create --port 6379:6379:tcp --private demo-env redis-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS          PUBLIC  URL
redislb16ae6     redis-lb           demo-env              6378:6379:TCP  false

```

The following is a summary of the arguments passed in the above command:

- `loadbalancer create`: creates a new load balancer
- `--port 6379:6379/TCP`: instructs the load balancer to forward requests from port 6379 on the server to port 6379 in the Docker container using the TCP protocol
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

### Part 3: Create the Service

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

### Part 4: Check the Status of the Service

As in the first deployment, we can keep an eye on our service by using the `service get` command:

`l0 service get redis-svc`

Once the service has finished scaling, try looking at the service's logs to see the output that the Redis server creates:

`l0 service logs redis-svc`

Among some warnings and information not important to this exercise and a fun bit of ASCII art, you should see following like the following:

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
            "value": ""
        }
    ],
    ...
```

That `value` is what will point the Guestbook application towards the Redis server. It needs to be populated in the following format:

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

Copy that `URL` value, paste it into the `value` field of `Guestbook.Dockerrun.aws.json`, append `:6379` to it, and save the file. It should look something like the following:

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

### Part 6: Update the Service

Almost all the pieces are in place! Now we just need to apply the new deploy to the running service:

`l0 service update guestbook-svc guestbook-dpl:latest`

As the service moves through the phases of its update process, we should see outputs like the following (if we keep an eye on the service with `l0 service get guestbook-svc`, that is):

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2*  1/1
                                                        guestbook-dpl:1
```
- _`guestbook-dpl:2` is in a transitional state_
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS      SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2  2/1
                                                        guestbook-dpl:1
```
- _both versions of the deployment are running at scale_
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2   1/1
                                                        guestbook-dpl:1*
```
- _`guestbook-dpl:1` is in a transitional state_
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS      SCALE
guestbo5fadd  guestbook-svc  demo-env     guestbook-lb  guestbook-dpl:2  1/1
```
- _`guestbook-dpl:1` has been removed, and only `guestbook-dpl:2` remains_


---

### Part 7: Prove It

You should now be able to point your browser at the URL for the Guestbook loadbalancer (run `l0 loadbalancer get guestbook-svc` to find it) and see what looks like the same Guestbook application you deployed in the first section of the walkthrough. Go ahead and add a few entries, make sure it's functioning properly. We'll wait.

Now, let's prove that we've actually divorced the data from the application by deleting and redeploying the Guestbook application:

`l0 service delete guestbook-svc`
`l0 loadbalancer delete guestbook-lb`

_(We'll leave the `deploy` intact so we can spin up a new service easily, and we'll leave the environment untouched because it also contained the Redis server.)_

Once those resources have been deleted (you can check the status of your Layer0 jobs with `l0 job list`), let's recreate them.

Create another load balancer:

`l0 loadbalancer create --ports 80:80/http demo-env guestbook-lb`

And another service, using the **guestbook-dpl** deploy we kept around:

`l0 service create --loadbalancer demo-env:guestbook-lb demo-env guestbook-svc guestbook-dpl:latest`

Wait for everything to spin up, and hit that new load balancer's url (`l0 loadbalancer get guestbook-lb`) with your browser. Your data should still be there!


---

### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Layer0 to delete the environment and terminate the application.

`l0 environment delete demo-env`

However, if you intend to continue through [Deployment 3](#deployment-3-guestbook-redis-consul), you will want to keep the resources you made in this section.


---

## 2b: Deploy with Terraform


---

### Part 1:


---

## Deployment 3: Guestbook + Redis + Consul

INTRO TEXT GOES HERE. You can choose to complete this section using either the [Layer0 CLI](#3a-deploy-with-layer0-cli) or [Terraform](#3b-deploy-with-terraform).


## 3a: Deploy with Layer0 CLI


---

### Part 1:


---

## 3b: Deploy with Terraform


---

### Part 1:


---

