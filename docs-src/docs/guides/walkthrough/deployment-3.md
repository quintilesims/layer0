# Deployment 3: Guestbook + Redis + Consul

In [Deployment 2](deployment-2#deploy-with-layer0-cli), we created two services in the same environment and linked them together manually.
While that can work for a small system, it's not really feasible for a system with a lot of moving parts - we would need to look up load balancer endpoints for all of our services and manually link them all together.
To that end, here we're going to to redeploy our two-service system using [Consul](https://www.consul.io), a service discovery tool.

For this deployment, we'll create a cluster of Consul servers which will be responsible for keeping track of the state of our system.
We'll also deploy new versions of the Guestbook and Redis task definition files - in addition to creating a container for its respective application, each task definition creates two other containers:

 - a container for a Consul agent, which is in charge of communicating with the Consul server cluster
 - a container for [Registrator](https://github.com/gliderlabs/registrator), which is charge of talking to the local Consul agent when a service comes up or goes down.

You can choose to complete this section using either the [Layer0 CLI](#deploy-with-layer0-cli) or [Terraform](#deploy-with-terraform).


## Deploy with Layer0 CLI

If you're following along, you'll want to be working in the `walkthrough/deployment-3/` directory of your clone of the [guides](https://github.com/quintilesims/guides) repo.

Files used in this deployment:

| Filename | Purpose |
|----------|---------|
| `CLI.Consul.Dockerrun.aws.json` | Template for running a Consul server |
| `CLI.Guestbook.Dockerrun.aws.json` | Template for running the Guestbook application with Registrator and Consul agent |
| `CLI.Redis.Dockerrun.aws.json` | Template for running a Redis server with Registrator and Consul agent |


---

### Part 1: Create the Consul Load Balancer

The Consul server cluster will live in the same environment as our Guestbook and Redis services - if you've completed [Deployment 1](deployment-1) and [Deployment 2](deployment-2), this environment already exists as **demo-env**.
We'll start by creating a load balancer for the Consul cluster.
The load balancer will be private since only Layer0 services need to communicate with the Consul cluster.
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
In `CLI.Consul.Dockerrun.aws.json`, find the entry in the `environment` block that looks like this:

```
{
    "name": "CONSUL_SERVER_URL",
    "value": ""
}
```

Update the "value" with the Consul load balancer's URL into and save the file.
We can then create the deploy.
At the command prompt, execute the following:

`l0 deploy create CLI.Consul.Dockerrun.aws.json consul-dpl`

We should see output like the following:

```
DEPLOY ID     DEPLOY NAME  VERSION
consul-dpl.1  consul-dpl   1
```

The following is a summary of the arguments passed in the above command:

- `deploy create`: creates a new Layer0 Deploy and allows you to specifiy an ECS task definition
- `CLI.Consul.Dockerrun.aws.json`: the file name of the ECS task definition (use the full path of the file if it is not in the current working directory)
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
In each of the CLI.Redis and CLI.Guestbook task definition files, look for the `CONSUL_SERVER_URL` block in the `consul-agent` container and populate the value field with the Consul load balancer's URL, then save the file.
At the command prompt, execute the two following commands to create new versions of the deploys for the Redis and Guestbook applications:

`l0 deploy create CLI.Redis.Dockerrun.aws.json redis-dpl`

`l0 deploy create CLI.Guestbook.Dockerrun.aws.json guestbook-dpl`

Then, execute the two following commands to redeploy the existing Redis and Guestbook services using those new deploys:

`l0 service update --wait redis-svc redis-dpl:latest`

`l0 service update --wait guestbook-svc guestbook-dpl:latest`

!!! NOTE
    Here, we should run `l0 service logs consul-svc` again and confirm that the Consul cluster has discovered these two services.

We can use `l0 loadbalancer get guestbook-lb` to obtain the guestbook application's URL, and then navigate to it with a web browser.
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

- `curl -s localhost:8500/v1/catalog/services`: use `curl` to send a GET request to the specified URL, where `localhost:8500` is an HTTP connection to the local Consul agent in this EC2 instance (the `-s` flag just silences excess output from `curl`)
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

To _really_ see how the Guestbook application connects to Redis, we can take an _even closer_ look!

Run `docker ps` to generate a listing of all the containers that Docker is running on the EC2 instance, and note the Container ID for the Guestbook container. Then run the following command to connect to the Guestbook container:

`docker exec -it [container_id] /bin/sh`

Once we've gotten inside the container, we'll run a similar command to the previous `curl`:

`curl -s consul-agent:8500/v1/catalog/service/redis`

Our Guestbook application makes a call like this one and figures out how to connect to the Redis service by mushing together the information from the `ServiceAddress` and `ServicePort` fields!

To close the `ssh` connection to the EC2 instance, run `exit` in the command prompt.


---

### Cleanup

When you're finished with the example, we can instruct Layer0 to terminate the applications and delete the environment.

`l0 environment delete demo-env`


---

## Deploy with Terraform

As before, we can complete this deployment using Terraform and the Layer0 provider instead of the Layer0 CLI.
As before, we will assume that you've cloned the [guides](https://github.com/quintilesims/guides) repo and are working in the `iterative-walkthrough/deployment-3/` directory.

We'll use these files to manage our deployment with Terraform:

| Filename | Purpose |
|----------|---------|
| `Guestbook.Dockerrun.aws.json` | Template for running the Guestbook application |
| `main.tf` | Provisions resources; populates variables in template files |
| `outputs.tf` | Values that Terraform will yield during deployment |
| `Redis.Dockerrun.aws.json` | Template for running the Redis application |
| `terraform.tfstate` | Tracks status of deployment _(created and managed by Terraform)_ |
| `terraform.tfvars` | Variables specific to the environment and application(s) |
| `variables.tf` | Values that Terraform will use during deployment |

---

### `*.tf`: A Brief Aside: Revisited: Redux

In looking at `main.tf`, you can see that we're pulling in a Consul module that we maintain (here's the [repo](https://github.com/quintilesims/consul)); this removes the need for a local task definition file.

We also are continuing to use modules for Redis and Guestbook.
However, instead of just sourcing the module and passing in a value or two, you can see that we actually create new deploys from local task definition files and pass those deploys in to the module.
This design allows us to use pre-made modules while also offering a great deal of flexibility.
If you'd like to follow along the Redis deployment logic chain (the other applications/services work similarly), it goes something like this:

- `main.tf` creates a deploy for the Redis server by rendering a local task definition and populating it with certain values
- `main.tf` passes the ID of the deploy into the Redis module, along with other values the module requires
- [the Redis module](https://github.com/quintilesims/redis/tree/master/terraform) pulls all the variables it knows about (both the defaults in `variables.tf` as well as the ones passed in)
- among other Layer0/AWS resources, the module spins up a Redis service; since a deploy ID has been provided, it uses that deploy to create the service instead of a deploy made from a [default task definition](https://github.com/quintilesims/redis/tree/master/terraform/Dockerrun.aws.json) contained within the module


---

### Part 1: Terraform Get

Run `terraform get` to pull down all the source materials Terraform needs for our deployment.


---


### Part 2: Terraform Init

This deployment has provider dependencies so an init call must be made. 
(Terraform v0.11~ requries init)
At the command prompt, execute the following command:

`terraform init`

We should see output like the following:

```
Initializing modules...
- module.consul
- module.redis
- module.guestbook

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

### Part 2: Terraform Plan

As before, we can run `terraform plan` to see what's going to happen.
We should see that there are 12 new resources to be created:

- the environment
- the two local deploys which will be used for Guestbook and Redis
- the load balancer, deploy, and service from each of the Consul, Guestbook, and Redis modules
    - _note that even though the default modules' deploys are created, they won't actually be used to deploy services_


---

### Part 3: Terraform Apply

Run `terraform apply`, and we should see output similar to the following:

```
data.template_file.consul: Refreshing state...
layer0_deploy.consul-dpl: Creating...

...
...
...

layer0_service.guestbook-svc: Creation complete

Apply complete! Resources: 10 added, 0 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: terraform.tfstate

Outputs:

guestbook_url = <http endpoint for the guestbook application>
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

