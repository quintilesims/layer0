# Deployment guide: Guestbook with a Database
The Guestbook application that you deployed in the [Guestbook deployment guide](/guides/guestbook) was a very simple application that stored its data in memory (also known as a "stateful" application). If you were to re-deploy the Guestbook service, all of the data previously entered into the application would be lost permanently.

A stateless application, on the other hand, does not record data generated in one session for use in subsequent sessions. In order to prevent data loss, web applications that you deploy using Layer0 should be stateless.

This guide will show you how to make a stateless Guestbook application that can communicate with an external resource (such as a database).

---

## Before you start
In order to complete the procedures in this section, you must install and configure Layer0. If you have not already configured Layer0, see the [installation guide](/setup/install).

[Install the Layer0 Terraform Plugin](/reference/terraform-plugin#install). The Layer0 Terraform Plugin makes Layer0 deployment information (like VPCs and subnets) available to Terraform configurations (.tf files).

---

## Deploy with Terraform
Use the Layer0 Terraform Plugin.


### Part 1: Download the configuration files
* [Dockerrun.aws.json](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/Dockerrun.aws.json)
* [terraform.tfvars](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/terraform.tfvars)
* [layer0.tf](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/layer0.tf)

!!! Note "External Resources"
	**This particular `layer0.tf` file uses a preconfigured DynamoDB resource as an example. Note that you can substitute any sort of external resource in its place, but instructions for configuring such an external resource are beyond the scope of this document. You can inspect the [relevant section](https://github.com/quintilesims/layer0-examples/tree/master/guestbook-db) of the layer0-examples repository to get a better understanding of what's going on in this guide.**


### Part 2: Terraform Apply
Run `terraform apply` to begin the process. Terraform will prompt you for configuration values that it does not have.

To begin deploying the application, run the following command:
<ul>
  <li class="command">`terraform apply`</li>
</ul>

_To avoid entering these values manually each time you run terraform, you can set the terraform variables by editing the `terraform.tfvars` file._

```
var.access_key
  Enter a value: [your AWS access key]

var.endpoint
  Enter a value: [your Layer0 endpoint]

var.secret_key
  Enter a value: [your AWS secret key]

var.token
  Enter a value: [your Layer0 token]

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
Terraform provisions the AWS resources (a DynamoDB instance, VPC and subnet configurations to connect the RDS instance to the Layer0 application), configures environment variables for the application, and deploys the application into a Layer0 environment.

You can use Terraform with Layer0 and AWS to create "fire and forget" deployments for your applications.

We use these files to set up a Layer0 envrionment with Terraform.

|Filename|Purpose|
|----|----|
|`terraform.tfvars`|Variables specific to the environment and guestbook application|
|`Dockerrun.aws.json`|Template for running the guestbook application in a Layer0 environment|
|`layer0.tf`|Provision Layer0 resources and populate variables in `Dockerrun.aws.json`|

Terraform figures out the appropriate order for creating each resource and handles the entire provisioning process.


### Cleanup
When you're finished with the example run `terraform destroy` in the same directory to destroy the AWS resources, Layer0 environment, and application.
