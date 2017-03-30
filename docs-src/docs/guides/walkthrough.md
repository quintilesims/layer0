# An Iterative Walkthrough

This guide aims to take you through three increasingly-complex deployment examples using Layer0. Successive sections build upon the previous ones, and each deployment can be completed either through using the Layer0 CLI directly, or through Terraform using our custom Layer0 Terraform Provider.

We assume that you're using Layer0 v0.9.0 or later. If you have not already installed and configured Layer0, see the [installation guide](/setup/install). If you are running an older version of Layer0, you may need to [upgrade](/setup/upgrade#upgrading-older-versions-of-layer0).

If you intend to deploy services using the Layer0 Terraform Provider, you'll want to make sure that you've [installed](/reference/terraform-plugin/#install) the provider correctly.

**Table of Contents**:

- [Deployment 1](#deployment-1-a-simple-guestbook-app): A Simple Guestbook App
- [Deployment 2](#deployment-2-guestbook-redis): Guestbook + Redis
- [Deployment 3](#deployment-3-guestbook-redis-consul): Guestbook + Redis + Consul


## Deployment 1: A Simple Guestbook App

In this section you'll learn how different Layer0 commands work together to deploy web applications to the cloud. The example application in this section is a guestbook -- a web application that acts as a simple message board. You can choose to complete this section using either the [Layer0 CLI](#1a-deploy-with-layer0-cli) or [Terraform](#1b-deploy-with-terraform).


## 1a: Deploy with Layer0 CLI

You will need to download the [Guestbook Task Definition](https://github.com/quintilesims/layer0-examples/blob/guides-overhaul/iterative-walkthrough/guestbook/Guestbook.Dockerrun.aws.json) file. Save that file as **Guestbook.Dockerrun.aws.json**. This walkthrough will assume that the file is located in your current working directory; if you choose to place the file elsewhere, you will need to provide the path to the file instead of just the filename when we reference it later in the guide.

---

### Part 1: Create the Environment

The first step in deploying an application with Layer0 is to create an environment. An environment is a dedicated space in which one or more services can reside. Here, we'll create a new environment named **demo-env**. At the command prompt, execute the following:

`l0 environment create demo-env`

This will give us the following output:

```
ENVIRONMENT ID  ENVIRONMENT NAME  CLUSTER COUNT  INSTANCE SIZE
demo00e6aa9     demo-env          0              m3.medium
```

---

### Part 2: Create the Load Balancer

In order to expose a web application to the public internet, we need to create a load balancer. A load balancer listens for web traffic at a specific address and directs that traffic to a Layer0 service.

A load balancer also has a notion of a health check - a way to assess whether or not the service is healthy and running properly. By default, Layer0 configures the health check of a load balancer based upon a simple TCP ping to port 80 every thirty seconds. Also by default, this ping will timeout after five seconds of no response from the service, and two consecutive successes or failures are required for the service to be considered healthy or unhealthy.

Here, we'll create a new load balancer named **guestbook-lb** inside of our environment named **demo-env**. The load balancer will listen on port 80, and forward that traffic along to port 80 in the Docker container using the HTTP protocol. Since the port configuration is already aligned with the default health check, we don't need to specify any health check configuration when we create this load balancer. At the command prompt, execute the following:

`l0 loadbalancer create --port 80:80/http demo-env guestbook-lb`

This will give us the following output:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE  PORTS       PUBLIC  URL
guestbodb65a     guestbook-lb       demo-env              80:80/HTTP  true
```

The following is a summary of the arguments passed in the above command:

* `loadbalancer create`: creates a new load balancer
* `--port 80:80/HTTP`: instructs the load balancer to forward requests from port 80 on the server to port 80 in the Docker container using the HTTP protocol
* `demo-env`: the name of the environment in which you are creating the load balancer
* `guestbook-lb`: a name for the load balancer itself

---

### Part 3: Deploy the Docker Task Definition

The `deploy` command is used to specify the Docker task definition that refers to a web application.

Here, we'll create a new deploy called **guestbook-dep** that refers to the **Guestbook.Dockerrun.aws.json** file you obtained earlier. At the command prompt, execute the following:

`l0 deploy create Guestbook.Dockerrun.aws.json guestbook-dep`

This will give us the following output:

```
DEPLOY ID        DEPLOY NAME    VERSION
guestbook-dep.1  guestbook-dep  1
```

The following is a summary of the arguments passed in the above command:

* `deploy create`: creates a new deployment and allows you to specify a Docker task definition
* `Guestbook.Dockerrun.aws.json`: the file name of the Docker task definition (use the full path of the file if it is not in your current working directory)
* `guestbook-dep`: a name for the deploy, which you will use later when you create the service

_The Deploy Name and Version are combined to create a unique identifier for a deploy. If you create additional deploys named **guestbook-dep**, they will be assigned different version numbers._


---

### Part 4: Create the Service

The final stage of the deployment process involves using the `service` command to create a new service and associate it with the environment, load balancer, and deploy that we created in the previous sections. The service will execute the Docker containers which have been described in the deploy.

Here, we'll create a new service called **guestbook-svc**. At the command prompt, execute the following:

`l0 service create --loadbalancer demo-env:guestbook-lb demo-env guestbook-svc guestbook-dep:latest`

This will give us the following output:

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dep:1*  0/1
```

The following is a summary of the arguments passed in the above command:

* `service create`: creates a new service
* `--loadbalancer demo-env:guestbook-lb`: the fully-qualified name of the load balancer; in this case, the load balancer named **guestbook-lb** in the environment named **demo-env**. 
	- _(It is not strictly necessary to use the fully qualified name of the load balancer, unless another load balancer with exactly the same name exists in a different environment.)_
* `demo-env`: the name of the environment you created in Part 1
* `guestbook-svc`: a name for the service you are creating
* `guestbook-dep`: the name of the deploy that you created in Part 3


---

### Check the Status of the Service

After a service has been created, it may take several minutes for that service to completely finish deploying. A service's status may be checked by using the `service get` command.

Let's take a peek at our **guestbook-svc** service. At the command prompt, execute the following:

`l0 service get demo-env:guestbook-svc`

If we're quick enough, we'll be able to see the first stage of the process (this is what was output after running the `service create` command up in Part 4). We should see an asterisk (\*) next to the name of the **guestbook-dep:1** deploy, which indicates that the service is in a transitional state:
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dep:1*  0/1
```

In the next phase of deployment, if we execute the `service get` command again, we will see **(1)** in the **Scale** column; this indicates that 1 copy of the service is transitioning to an active state:
```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dep:1*  0/1 (1)
```

In the final phase of deployment, we should see that the service has been fully scaled up to the desired count. If we execute the `service get` command again, you will see the following output:

```
SERVICE ID    SERVICE NAME   ENVIRONMENT  LOADBALANCER  DEPLOYMENTS       SCALE
guestbo9364b  guestbook-svc  demo-env     guestbook-lb  guestbook-dep:1   1/1
```


---

### Get the Application's URL

Once the service has been completely deployed, we can obtain the URL for the application and launch it in a browser.

At the command prompt, execute the following:

`l0 loadbalancer get demo-env:guestbook-lb`

This will give us the following output:

```
LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICE        PORTS       PUBLIC  URL
guestbodb65a     guestbook-lb       demo-env     guestbook-svc  80:80/HTTP  true    <url>
```

Copy the value shown in the **URL** column and paste it into a web browser. The guestbook application will appear (once the service has completely finished deploying).


---

### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Layer0 to delete the environment and terminate the application and the Layer0 the application was using.

`l0 environment delete demo-env`

However, if you intend to continue through [Deployment 2](#deployment-2-guestbook-redis), you will want to keep the resources you made in this section.


---

## 1b: Deploy with Terraform

Instead of using the Layer0 CLI directly, you can instead use the Layer0 Provider for Terraform, and deploy using Terraform. 


---

### Part 1: Download the Configuration Files

- [Guestbook.Dockerrun.aws.json](https://github.com/quintilesims/layer0-examples/blob/guides-overhaul/iterative-walkthrough/guestbook/Guestbook.Dockerrun.aws.json)
- [terraform.tfvars](https://github.com/quintilesims/layer0-examples/blob/guides-overhaul/iterative-walkthrough/terraform/terraform.tfvars)
- [layer0.tf](https://github.com/quintilesims/layer0-examples/blob/guides-overhaul/iterative-walkthrough/terraform/deployment-1/layer0.tf) for Deployment 1


---

### Part 2: Terraform Apply

Run `terraform apply` to begin the process. Terraform will prompt you for configuration values that it does not have.

To begin deploying the application, run the following command:

`terraform apply`

_To avoid entering these values manually each time you run terraform, you can set the terraform variables by editing the `terraform.tfvars` file._

```
var.endpoint
  Enter a value: <enter your Layer0 endpoint>

var.token
  Enter a value: <enter your Layer0 token>

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

__It may take a few minutes for the guestbook service to launch and the load balancer to become available. During that time you may get HTTP 503 errors when making HTTP requests against the load balancer URL.__

Terraform will set up the entire environment for you and then output a link to the application's load balancer.


### What's happening

Terraform provisions the AWS resources (an RDS instance, VPC and subnet configurations to connect the RDS instance to the Layer0 application), configures environment variables for the application, and deploys the application into a Layer0 environment.

You can use Terraform with Layer0 and AWS to create "fire-and-forget" deployments for your applications.

We use these files to set up a Layer0 envrionment with Terraform:

|Filename|Purpose|
|----|----|
|`terraform.tfvars`|Variables specific to the environment and guestbook application|
|`Guestbook.Dockerrun.aws.json`|Template for running the guestbook application in a Layer0 environment|
|`layer0.tf`|Provision Layer0 resources and populate variables in `Guestbook.Dockerrun.aws.json`|

Terraform figures out the appropriate order for creating each resource and handles the entire provisioning process.


### Cleanup

If you're finished with the example and don't want to continue with this walkthrough, you can instruct Terraform to destroy the AWS resources, the Layer0 environment, and the application. Execute the following command (in the same directory):

`terraform destroy`

However, if you intend to continue through [Deployment 2](#deployment-2-guestbook-redis), you will want to keep the resources you made in this section.


---

## Deployment 2: Guestbook + Redis

INTRO TEXT GOES HERE. You can choose to complete this section using either the [Layer0 CLI](#2a-deploy-with-layer0-cli) or [Terraform](#2b-deploy-with-terraform).


## 2a: Deploy with Layer0 CLI


---

### Part 1:


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

