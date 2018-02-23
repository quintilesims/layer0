# Deployment 2: Guestbook + Redis

In this section, we're going to add some complexity to the previous deployment.
[Deployment 1](deployment-1) saw us create a simple guestbook application which kept its data in memory.
But what if that ever came down, either by intention or accident?
It would be easy enough to redeploy it, but all of the entered data would be lost.
What if we wanted to scale the application to run more than one copy?
For this deployment, we're going to separate the data store from the guestbook application by creating a second Layer0 service which will house a Redis database server and linking it to the first.
You can choose to complete this section using either [the Layer0 CLI](#deploy-with-layer0-cli) or [Terraform](#deploy-with-terraform).


---

## Deploy with Layer0 CLI

For this example, we'll be working in the `walkthrough/deployment-2/` directory of the [guides](https://github.com/quintilesims/guides/) repo.
We assume that you've completed the [Layer0 CLI](deployment-1#deploy-with-layer0-cli) section of Deployment 1.

Files used in this deployment:

| Filename | Purpose |
|----------|---------|
| `Guestbook.Dockerrun.aws.json` | Template for running the Guestbook application |
| `Redis.Dockerrun.aws.json` | Template for running a Redis server |


---

### Part 1: Create the Redis Load Balancer

Both the Guestbook service and the Redis service will live in the same Layer0 environment, so we don't need to create one like we did in the first deployment.
We'll start by making a load balancer behind which the Redis service will be deployed.

The `Redis.Dockerrun.aws.json` task definition file we'll use is very simple - it just spins up a Redis server with the default configuration, which means that it will be serving on port 6379.
Our load balancer needs to be able to forward TCP traffic to and from this port.
And since we don't want the Redis server to be exposed to the public internet, we'll put it behind a private load balancer; private load balancers only accept traffic that originates from within their own environment.
We'll also need to specify a non-default healthcheck target, since the load balancer won't expose port 80.
At the command prompt, execute the following:

`l0 loadbalancer create --port 6379:6379/tcp --private --healthcheck-target tcp:6379 demo-env redis-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS          PUBLIC  URL
redislb16ae6     redis-lb           demo-env              6378:6379:TCP  false
```

The following is a summary of the arguments passed in the above command:

- `loadbalancer create`: creates a new load balancer
- `--port 6379:6379/TCP`: instructs the load balancer to forward requests from port 6379 on the load balancer to port 6379 in the EC2 instance using the TCP protocol
- `--private`: instructs the load balancer to ignore external traffic
- `--healthcheck-target tcp:6379`: instructs the load balancer to check the health of the service via TCP pings to port 6379
- `demo-env`: the name of the environment in which the load balancer is being created
- `redis-lb`: a name for the load balancer itself


---

### Part 2: Deploy the ECS Task Definition

Here, we just need to create the deploy using the `Redis.Dockerrun.aws.json` task definition file.
At the command prompt, execute the following:

`l0 deploy create Redis.Dockerrun.aws.json redis-dpl`

We should see output like the following:

```
DEPLOY ID    DEPLOY NAME  VERSION
redis-dpl.1  redis-dpl    1
```

The following is a summary of the arguments passed in the above command:

- `deploy create`: creates a new Layer0 Deploy and allows you to specify an ECS task definition
- `Redis.Dockerrun.aws.json`: the file name of the ECS task definition (use the full path of the file if it is not in your current working directory)
- `redis-dpl`: a name for the deploy, which we will use later when we create the service


---

### Part 3: Create the Redis Service

Here, we just need to pull the previous resources together into a service.
At the command prompt, execute the following:

`l0 service create --wait --loadbalancer demo-env:redis-lb demo-env redis-svc redis-dpl:latest`

We should see output like the following:

```
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
redislb16ae6  redis-svc     demo-env     redis-lb      redis-dpl:1  0/1
```

The following is a summary of the arguments passed in the above commands:

- `service create`: creates a new Layer0 Service
- `--wait`:  instructs the CLI to keep hold of the shell until the service has been successfully deployed
- `--loadbalancer demo-env:redis-lb`: the fully-qualified name of the load balancer; in this case, the load balancer named **redis-lb** in the environment named **demo-env**
    - _(Again, it's not strictly necessary to use the fully-qualified name of the load balancer as long as there isn't another load balancer with the same name in a different environment)_
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

You should see in `walkthrough/deployment-2/` another `Guestbook.Dockerrun.aws.json` file.
This file is very similar to but not the same as the one in `deployment-1/` - if you open it up, you can see the following additions:

```
    ...
    "environment": [
        {
            "name": "GUESTBOOK_BACKEND_TYPE",
            "value": "redis"
        },
        {
            "name": "GUESTBOOK_BACKEND_CONFIG",
            "value": "<redis host and port here>"
        }
    ],
    ...
```

The `"GUESTBOOK_BACKEND_CONFIG"` variable is what will point the Guestbook application towards the Redis server.
The `<redis host and port here>` section needs to be replaced and populated in the following format:

```
"value": "ADDRESS_OF_REDIS_SERVER:PORT_THE_SERVER_IS_SERVING_ON"
```

We already know that Redis is serving on port 6379, so let's go find the server's address.
Remember, it lives behind a load balancer that we made, so run the following command:

`l0 loadbalancer get redis-lb`

We should see output like the following:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE    PORTS          PUBLIC  URL
redislb16ae6     redis-lb           demo-env     redis-svc  6379:6379/TCP  false   internal-l0-<yadda-yadda>.elb.amazonaws.com
```

Copy that `URL` value, replace `<redis host and port here>` with the `URL` value in `Guestbook.Dockerrun.aws.json`, append `:6379` to it, and save the file.
It should look something like the following:

```
    ...
    "environment": [
        {
            "name": "GUESTBOOK_BACKEND_CONFIG",
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

Almost all the pieces are in place!
Now we just need to apply the new Guestbook deploy to the running Guestbook service:

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

You should now be able to point your browser at the URL for the Guestbook load balancer (run `l0 loadbalancer get guestbook-lb` to find it) and see what looks like the same Guestbook application you deployed in the first section of the walkthrough.
Go ahead and add a few entries, make sure it's functioning properly.
We'll wait.

Now, let's prove that we've actually separated the data from the application by deleting and redeploying the Guestbook application:

`l0 service delete --wait guestbook-svc`

_(We'll leave the `deploy` intact so we can spin up a new service easily, and we'll leave the environment untouched because it also contained the Redis server.
We'll also pass the `--wait` flag so that we don't need to keep checking on the status of the job to know when it's complete.)_

Once those resources have been deleted, we can recreate them!

Create another service, using the **guestbook-dpl** deploy we kept around:

`l0 service create --loadbalancer demo-env:guestbook-lb demo-env guestbook-svc guestbook-dpl:latest`

Wait for everything to spin up, and hit that new load balancer's url (`l0 loadbalancer get guestbook-lb`) with your browser.
Your data should still be there!


---

### Cleanup

When you're finished with the example, you can instruct Layer0 to delete the environment and terminate the application.

`l0 environment delete demo-env`


---

## Deploy with Terraform

As before, we can complete this deployment using Terraform and the Layer0 provider instead of the Layer0 CLI. As before, we will assume that you've cloned the [guides](https://github.com/quintilesims/guides) repo and are working in the `walkthrough/deployment-2/` directory.

We'll use these files to manage our deployment with Terraform:

| Filename | Purpose |
|----------|---------|
| `main.tf` | Provisions resources; populates variables in template files |
| `outputs.tf` | Values that Terraform will yield during deployment |
| `terraform.tfstate` | Tracks status of deployment _(created and managed by Terraform)_ |
| `terraform.tfvars` | Variables specific to the environment and application(s) |
| `variables.tf` | Values that Terraform will use during deployment |


---

### `*.tf`: A Brief Aside: Revisited

Not much is changed from [Deployment 1](deployment-1#deploy-with-terraform).
In `main.tf`, we pull in a new, second module that will deploy Redis for us.
We maintain this module as well; you can inspect [the repo](https://github.com/quintilesims/redis) if you'd like.

In `main.tf` where we pull in the Guestbook module, you'll see that we're supplying more values than we did last time, because we need some additional configuration to let the Guestbook application use a Redis backend instead of its default in-memory storage.


---

### Part 1: Terraform Get

Run `terraform get` to pull down the source materials Terraform will use for deployment.
This will create a local `.terraform/` directory.


---


### Part 2: Terraform Init

This deployment has provider dependencies so an init call must be made. 
(Terraform v0.11~ requries init)
At the command prompt, execute the following command:

`terraform init`

We should see output like the following:

```
Initializing modules...
- module.redis
  Getting source "github.com/quintilesims/redis//terraform"
- module.guestbook
  Getting source "github.com/quintilesims/guides//guestbook/module"

Initializing provider plugins...
- Checking for available provider plugins on https://releases.hashicorp.com...
- Downloading plugin for provider "template" (1.0.0)...

The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.

* provider.template: version = "~> 1.0"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```
---

### Part 3: Terraform Plan

It's always a good idea to find out what Terraform intends to do, so let's do that:

`terraform plan`

As before, we'll be prompted for any variables Terraform needs and doesn't have (see the note in [Deployment 1](deployment-1#deploy-with-terraform) for configuring Terraform variables).
We'll see output similar to the following:

```
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.template_file.redis: Refreshing state...
The Terraform execution plan has been generated and is shown below.
Resources are shown in alphabetical order for quick scanning. Green resources
will be created (or destroyed and then created if an existing resource
exists), yellow resources are being changed in-place, and red resources
will be destroyed. Cyan entries are data sources to be read.

Note: You didn't specify an "-out" parameter to save this plan, so when
"apply" is called, Terraform can't guarantee this is what will execute.

+ layer0_environment.demo
    ami:               "<computed>"
    current_scale:     "<computed>"
    name:              "demo"
    os:                "linux"
    security_group_id: "<computed>"
    size:              "m3.medium"

+ module.redis.layer0_deploy.redis
    content: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"redis\",\n            \"image\": \"redis:3.2-alpine\",\n            \"essential\": true,\n            \"memory\": 128,\n            \"portMappings\": [\n                {\n                    \"hostPort\": 6379,\n              \"containerPort\": 6379\n                }\n            ]\n        }\n    ]\n}\n\n"
    name:    "redis"

+ module.redis.layer0_load_balancer.redis
    environment:                    "${var.environment_id}"
    health_check.#:                 "<computed>"
    name:                           "redis"
    port.#:                         "1"
    port.1072619732.certificate:    ""
    port.1072619732.container_port: "6379"
    port.1072619732.host_port:      "6379"
    port.1072619732.protocol:       "tcp"
    private:                        "true"
    url:                            "<computed>"

+ module.redis.layer0_service.redis
    deploy:        "${ var.deploy_id == \"\" ? layer0_deploy.redis.id : var.deploy_id }"
    environment:   "${var.environment_id}"
    load_balancer: "${layer0_load_balancer.redis.id}"
    name:          "redis"
    scale:         "1"
    wait:          "true"

<= module.guestbook.data.template_file.guestbook
    rendered: "<computed>"
    template: "{\n    \"AWSEBDockerrunVersion\": 2,\n    \"containerDefinitions\": [\n        {\n            \"name\": \"guestbook\",\n            \"image\": \"quintilesims/guestbook\",\n            \"essential\": true,\n       \"memory\": 128,\n            \"environment\": [\n                {\n                    \"name\": \"GUESTBOOK_BACKEND_TYPE\",\n                    \"value\": \"${backend_type}\"\n                },\n                {\n                    \"name\": \"GUESTBOOK_BACKEND_CONFIG\",\n                    \"value\": \"${backend_config}\"\n                },\n                {\n                    \"name\": \"AWS_ACCESS_KEY_ID\",\n  \"value\": \"${access_key}\"\n                },\n                {\n                    \"name\": \"AWS_SECRET_ACCESS_KEY\",\n                    \"value\": \"${secret_key}\"\n                },\n                {\n            \"name\": \"AWS_REGION\",\n                    \"value\": \"${region}\"\n                }\n   ],\n            \"portMappings\": [\n                {\n                    \"hostPort\": 80,\n     \"containerPort\": 80\n                }\n            ]\n        }\n    ]\n}\n"
    vars.%:   "<computed>"

+ module.guestbook.layer0_deploy.guestbook
    content: "${data.template_file.guestbook.rendered}"
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
    scale:         "2"
    wait:          "true"


Plan: 7 to add, 0 to change, 0 to destroy.
```

We should see that Terraform intends to add 7 new resources, some of which are for the Guestbook deployment and some of which are for the Redis deployment.


---

### Part 4: Terraform Apply

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
    It may take a few minutes for the guestbook service to launch and the load balancer to become available.
    During that time you may get HTTP 503 errors when making HTTP requests against the load balancer URL.


### What's Happening

Terraform provisions the AWS resources through Layer0, configures environment variables for the application, and deploys the application into a Layer0 environment.
Terraform also writes the state of your deployment to the `terraform.tfstate` file (creating a new one if it's not already there).


### Cleanup

When you're finished with the example, you can instruct Terraform to destroy the Layer0 environment, and terminate the application.
Execute the following command (in the same directory):

`terraform destroy`

It's also now safe to remove the `.terraform/` directory and the `*.tfstate*` files.


---

