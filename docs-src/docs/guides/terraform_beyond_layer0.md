# Deployment guide: Terraform beyond Layer0

 In this example, we'll learn how you can use Terraform to create a Layer0 service as well as a persistent data store. The main goal of this example is to explore how you can combine Layer0 with other Terraform providers and best practices.

## Before you start

To complete the procedures in this section, you must have the following installed and configured correctly:

 * Layer0 v0.8.4 or later
 * Terraform v0.9.0 or later
 * Layer0 Terraform Provider

If you have not already configured Layer0, see the [Layer0 installation guide](/setup/install). If you are running an older version of Layer0, see the [Layer0 upgrade instructions](/setup/upgrade#upgrading-older-versions-of-layer0).

See the [Terraform installation guide](/reference/terraform-plugin#install) to install Terraform and the Layer0 Terraform Plugin.

---

## Deploy with Terraform

Using Terraform, you will deploy a simple guestbook application backed by AWS DynamoDB Table. The terraform configuration file will use both the Layer0 and AWS Terraform providers, to deploy the guestbook application and provision a new DynamoDB Table.

## Part 1: Clone the guides repository

Run this command to clone the `quintilesims/guides` repository:

`git clone https://github.com/quintilesims/guides.git`

Once you have cloned the repository, navigate to the `guides/terraform-beyond-layer0/example-1` folder for the rest of this example.


## Part 2: Terraform Plan

!!! Note
	As we're using modules in our Terraform configuration, we need to run `terraform get` command before performing other terraform operations. Running `terraform get` will download the modules to your local folder named `.terraform`. See here for more information on [terraform get](https://www.terraform.io/docs/commands/get.html).

	`terraform get`

	Get: file:///Users/<username>/go/src/github.com/quintilesims/guides/terraform-beyond-layer0/example-1/modules/guestbook_service

Before deploying, we can run the following command to see what changes Terraform will make to your infrastructure should you go ahead and apply. If you had any errors in your layer0.tf file, running `terraform plan` would output those errors so that you can address them. Also, Terraform will prompt you for configuration values that it does not have.

!!! Tip
	There are a few ways to configure Terraform so that you don't have to keep entering these values every time you run a Terraform command (editing the `terraform.tfvars` file, or exporting environment variables like `TF_VAR_endpoint` and `TF_VAR_token`, for example). See the [Terraform Docs](https://www.terraform.io/docs/configuration/variables.html) for more.

`terraform plan`

```
var.endpoint
  Enter a value: <enter your Layer0 endpoint>

var.token
  Enter a value: <enter your Layer0 token>
...
+ aws_DynamoDB_table.guestbook
    arn:                       "<computed>"
    attribute.#:               "1"
    attribute.4228504427.name: "id"
    attribute.4228504427.type: "S"
    hash_key:                  "id"
    name:                      "guestbook"
    read_capacity:             "20"
    stream_arn:                "<computed>"
    stream_enabled:            "<computed>"
    stream_view_type:          "<computed>"
    write_capacity:            "20"

...
```

## Part 3: Terraform Apply

Run the following command to begin the deploy process.

`terraform apply`

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
	It may take a few minutes for the guestbook service to launch and the load balancer to become available. During that time, you may get HTTP 503 errors when making HTTP requests against the load balancer URL.


Terraform will set up the entire environment for you and then output a link to the application's load balancer.

### What's happening

Terraform using the [AWS provider](https://www.terraform.io/docs/providers/aws/index.html), provisions a new DynamoDB table. It also uses the [Layer0 provider](http://layer0.ims.io/reference/terraform-plugin/#provider) to provision the environment, deploy, load balancer and service required to run the entire guestbook application.

Looking at an excerpt of the file [./terraform-beyond-layer0/example-1/modules/guestbook_service/main.tf](https://github.com/quintilesims/guides/blob/master/terraform-beyond-layer0/example-1/modules/guestbook_service/main.tf), we can see the following definitions:

```
resource "aws_dynamodb_table" "guestbook" {
  name           = "${var.table_name}"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "layer0_deploy" "guestbook" {
  name    = "guestbook"
  content = "${data.template_file.guestbook.rendered}"
}

data "template_file" "guestbook" {
  template = "${file("Dockerrun.aws.json")}"

  vars {
    access_key = "${var.access_key}"
    secret_key = "${var.secret_key}"
    region     = "${var.region}"
    table_name = "${aws_dynamodb_table.guestbook.name}"
  }
}
```

Note the resource definitions for `aws_dynamodb_table` and `layer0_deploy`. To configure the guestbook application to use the provisioned DynamoDB table, we reference the `name` property from the DynamoDB definition `table_name = "${aws_dynamodb_table.guestbook.name}"`. 

These `vars` are used to populate the template fields in our [Dockerrun.aws.json](https://github.com/quintilesims/guides/blob/master/terraform-beyond-layer0/example-1/modules/guestbook_service/Dockerrun.aws.json) file. 

```
{
    "AWSEBDockerrunVersion": 2,
    "containerDefinitions": [
        {
            "name": "guestbook",
            "image": "quintilesims/guestbook-db",
            "essential": true,
            "memory": 128,
            "environment": [
                {
                    "name": "DYNAMO_TABLE",
                    "value": "${table_name}"
                }
                ...
```

The Layer0 configuration referencing the AWS DynamoDB configuration `table_name = "${aws_DynamoDB_table.guestbook.name}"`, infers an implicit dependency. Before Terraform creates the infrastructure, it will use this information to order the resource creation and create resources in parallel, where there are no dependencies. In this example, the AWS DynamoDB table will be created before the Layer0 deploy. See [Terraform Resource Dependencies](https://www.terraform.io/intro/getting-started/dependencies.html) for more information.

## Part 4: Scaling a Layer0 Service

The workflow to make changes to your infrastructure generally involves updating your Terraform configuration file followed by a `terraform plan` and `terraform apply`.

### Update the Terraform configuration
Open the file `./example-1/modules/guestbook_service/main.tf` in a text editor and make the change to add a `scale` property with a value of `3` to the `layer0_service` section. For more information about the `scale` property, see [Layer0 Terraform Plugin](http://layer0.ims.io/reference/terraform-plugin/#service) documentation. The result should look like the below:

example-1/modules/guestbook_service/main.tf

```
# Create a service named "guestbook"
resource "layer0_service" "guestbook" {
  name          = "guestbook"
  environment   = "${layer0_environment.demo.id}"
  deploy        = "${layer0_deploy.guestbook.id}"
  load_balancer = "${layer0_load_balancer.guestbook.id}"
  scale         = 3
}

```

### Plan and Apply

Execute the `terraform plan` command to understand the changes that you will be making. Note that if you did not specify `scale`, it defaults to '1'.

`terraform plan`

Outputs:

```
...

~ module.guestbook.layer0_service.guestbook
    scale: "1" => "3"
```

Now run the following command to deploy your changes:

`terraform apply`

Outputs:

```
layer0_environment.demo: Refreshing state... (ID: demoenvbb9f6)
data.template_file.guestbook: Refreshing state...
layer0_deploy.guestbook: Refreshing state... (ID: guestbook.6)
layer0_load_balancer.guestbook: Refreshing state... (ID: guestbo43ab0)
layer0_service.guestbook: Refreshing state... (ID: guestboebca1)
layer0_service.guestbook: Modifying... (ID: guestboebca1)
  scale: "1" => "3"
layer0_service.guestbook: Modifications complete (ID: guestboebca1)

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: 

Outputs:

services = <guestbook_service_url>
```

To confirm your service has been updated to the desired scale, you can run the following layer0 command. Note that the desired scale for the guestbook service should be eventually be 3/3.

`l0 service get guestbook1_guestbook_svc`
Outputs:

```
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
SERVICE ID    SERVICE NAME              ENVIRONMENT  LOADBALANCER             DEPLOYMENTS                  SCALE
guestbo4fd3b  guestbook1_guestbook_svc  demo         guestbook1_guestbook_lb  guestbook1_guestbook_dpl:3*  1/3 (2)
```

As scale is a parameter we are likely to change in the future, rather than hardcoding it to 3 as we have done just now, it would be better to use a variable to store  `service_scale`. The following Best Practices sections will show how you can achieve this.

!!! Note "Best Practices with Terraform + Layer0"
	The following sections outline some of the best practices and tips to take into consideration, when using Layer0 with Terraform.

## Part 5: Terraform Remote State

Terraform stores the state of the deployed infrastructure in a local file named `terraform.tfstate` by default. To find out more about why Terraform needs to store state, see [Purpose of Terraform State](https://www.terraform.io/docs/state/purpose.html). 

How state is loaded and used for operations such as `terraform apply` is determined by a [Backend](https://www.terraform.io/docs/backends). As mentioned, by default the state is stored locally which is enabled by a "local" backend.

### Remote State

By default, Terraform stores state locally but it can also be configured to store state in a remote backend. This can prove useful when you are working as part of a team to provision and manage services deployed by Terraform. All the members of the team will need access to the state file to apply new changes and be able to do so without overwriting each others' changes. See here for more information on the different [backend types](https://www.terraform.io/docs/backends/types/index.html) supported by Terraform.

To configure a remote backend, append the `terraform` section below to your terraform file `./example-1/main.tf`. Populate the `bucket` property to an existing s3 bucket.

!!! Tip
	If you have been following along with the guide, `./example-1/main.tf` should already have the below section commented out. You can uncomment the `terraform` section and populate the bucket property with an appropriate value.

```
terraform {
  backend "s3" {
    bucket     = "<my-bucket-name>"
    key        = "demo-env/remote-backend/terraform.tfstate"
    region     = "us-west-2"
  }
}
```

Once you have modified `main.tf`, you will need to initialize the newly configured backend by running the following command.

`terraform init`

Outputs:

```
Initializing the backend...

Do you want to copy state from "local" to "consul"?
  ...
  Do you want to copy the state from "local" to "consul"? Enter "yes" to copy
  and "no" to start with the existing state in "consul".

  Enter a value: 

```

Go ahead and enter: `yes`.

```

Successfully configured the backend "consul"! Terraform will automatically
use this backend unless the backend configuration changes.

Terraform has been successfully initialized!
...

```

### What's happening

As you are configuring a backend for the first time, Terraform will give you an option to migrate your state to the new backend. From now on, any further changes to your infrastructure made by Terraform will result in the remote state file being updated. For more information see [Terraform backends](https://www.terraform.io/docs/backends/index.html).

A new team member can use the `main.tf` from their own machine without obtaining a copy of the state file `terraform.tfstate` as the configuration will retrieve the state file from the remote backend.

### Locking

Not all remote backends support locking (locking ensures only one person is able to change the state at a time). The `S3` backend we used earlier in the example supports locking which is disabled by default. The `S3` backend uses a DynamoDB table to acquire a lock before making a change to the state file. To enable locking, you need to specify `locking_table` property with the name of an existing DynamoDB table. The DynamoDB table also needs primary key named `LockID` of type `String`.

### Security

A Terraform state file is written in plain text. This can lead to a situation where deploying resources that require sensitive data can result in the sensitive data being stored in the state file. To minimize exposure of sensitive data, you can enable [server side encryption](https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html) of the state file by adding property `encrypt` set to `true`.

This will ensure that the file is encrypted in S3 and by using a remote backend, you will also have the added benefit of the state file not being persisted to disk locally as it will only ever be held in memory by Terraform.

For securing the state file further, you can also enable access logging on the S3 bucket you are using for the remote backend, which can help track down invalid access should it occur.

## Part 6: Terraform Configuration Structure

While there are many different approaches to organizing your Terraform code, we suggest using the following file structure:

```
example1/  # contains overarching Terraform deployment, pulls in any modules that might exist
  ─ main.tf  
  ─ variables.tf  
  ─ output.tf  
  + modules/  # if you can break up deployment into smaller modules, keep the modules in here
      + guestbook_service/  # contains Terraform configuration for a module
        ─ main.tf  
        ─ variables.tf  
        ─ output.tf
      + service2/  # contains another module
      + service3/  # contains another module
```

Here we are making use of Terraform [Modules](https://www.terraform.io/docs/modules/index.html). Modules in Terraform are self-contained packages of Terraform configurations, that are managed as a group. Modules are used to create reusable components in Terraform as well as for basic code organization. In this example, we are using modules to separate each service and making it consumable as a module.

If you wanted to add a new service, you can create a new service folder inside the ./modules. If you wanted to you could even run multiple copies of the same service. See here for more information about [Creating Modules](https://www.terraform.io/docs/modules/create.html).

Also see the below repositories for ideas on different ways you can organize your Terraform configuration files for the needs of your specific project: 

* [Terraform Community Modules](https://github.com/terraform-community-modules)
* [Best Pratices Ops](https://github.com/hashicorp/best-practices)

## Part 7: State Environments

Layer0 recommends that you typically make a single environment for each tier of your application, such as `dev`, `staging` and `production`. That recommendation still holds when using Terraform with Layer0. Using Layer0 CLI, you can target a specific environment for most CLI commands. This enables you to service each tier relatively easily. In Terraform, there a few approaches you can take to enable a similar workflow.

### Single Terraform Configuration

You can use a single Terraform configuration to create and maintain multiple environments by making use of the [Count](https://www.terraform.io/docs/configuration/resources.html#count) parameter, inside a Resource. Count enables you to create multiple copies of a given resource. 

For example

```
variable "environments" {
  type = "list"

  default = [
    "dev",
    "staging"
    "production"
  ]
}

resource "layer0_environment" "demo" {
  count = "${length(var.environments)}"

  name = "${var.environments[count.index]}_demo"
}
```

Let's have a more in-depth look in how this works. You can start by navigating to `./terraform-beyond-layer0/example-2' folder. Start by running the plan command.

`terraform plan`

Outputs:
```
+ module.environment.aws_dynamodb_table.guestbook.0
    ...
    name:                      "dev_guestbook"
...
+ module.environment.aws_dynamodb_table.guestbook.1
    ..
    name:                      "staging_guestbook"
...
```

Note that you will see a copy of each resource for each environment specified in your environments file in `./example-2/variables.tf`. Go ahead and run apply.

`terraform apply`

Outputs:
```
Apply complete! Resources: 10 added, 0 changed, 0 destroyed.

Outputs:

guestbook_urls = 
<dev_url>
<staging_url>
```

You have now created two separate environments using a single terraform configuration: dev & staging. You can navigate to both the urls output and you should note that they are separate instances of the guestbook application backed with their own separate data store.

A common use case for maintaining different environments is to configure each environment slightly differently. For example, you might want to scale your Layer0 service to 3 for staging and leave it as 1 for the dev environment. This can be done easily by using conditional logic to set our `scale` parameter in the layer0 service configuration in `./example-2/main.tf`. Go ahead and open `main.tf` in a text editor. Navigate to the `layer0_service guestbook` section. Uncomment the scale parameter so that your configuration looks like below.

```
resource "layer0_service" "guestbook" {
  count = "${length(var.environments)}"

  name          = "${element(layer0_environment.demo.*.name, count.index)}_guestbook_svc"
  environment   = "${element(layer0_environment.demo.*.id, count.index)}"
  deploy        = "${element(layer0_deploy.guestbook.*.id, count.index)}"
  load_balancer = "${element(layer0_load_balancer.guestbook.*.id, count.index)}"
  scale         = scale         = "${lookup(var.service_scale, var.environments[count.index]), "1")}"
}
```

The variable `service_scale` is already defined in `variables.tf`. If you now go ahead and run plan, you will see that the `guestbook` service for only the `staging` environment will be scaled up.

`terraform plan`

Outputs:

```
~ layer0_service.guestbook.1
    scale: "1" => "3"
```

A potential downside of this approach however is that all your environments are using the same state file. Sharing a state file breaks some of the resource encapsulation between environments. Should there ever be a situation where your state file becomes corrupt, it would affect your ability to service all the environments till you resolve the issue by potentially rolling back to a previous copy of the state file. 

The next section will show you how you can separate your Terraform environment configuration such that each environment will have its own state file.

!!! Note
	As previously mentioned, you will want to avoid hardcoding resource parameter configuration values as much as possible. As an example the scale property of a layer0 service. But this extends to other properties as well like docker image version etc. You should avoid using `latest` and specify a explicit version via configurable variable when possible.

### Multiple Terraform Configurations

The previous example used a single set of Terraform Configuration files to create and maintain multiple environments. This resulted in a single state file which had the state information for all the environments. To avoid all environments sharing a single state file, you can split your Terraform configuration so that you a state file for each environment.

Go ahead and navigate to `./terraform-beyond-layer0/example-3` folder. Here we are using a folder to separate each environment. So `env-dev` and `env-staging` represent a `dev` and `staging` environment. To work with either of the environments, you will need to navigate into the desired environment's folder and run Terraform commands. This will ensure that each environment will have its own state file.

Open the env-dev folder inside a text editor. Note that `main.tf` doesn't contain any resource definitions. Instead, we only have one module definition which has various variables being passed in, which is also how we are passing in the `environment` variable. To create a `dev` and `staging` environments for our guestbook application, go ahead and run terraform plan and apply commands from `env-dev` and `env-staging` folders.

```
# assuming you are in the terraform-beyond-layer0/example-3 folder
cd env-dev
terraform get
terraform plan
terraform apply

cd ../env-staging
terraform get
terraform plan
terraform apply
```

You should now have two instances of the guestbook application running. Note that our guestbook service in our staging environment has been scaled to 3. We have done this by specifying a map variable `service_scale` in `./example-3/dev-staging/variables.tf` which can have different scale values for each environment.

## Part 8: Multiple Provider Instances

You can define multiple instances of the same provider that is uniquely customized. For example, you can have an `aws` provider to support multiple regions, different roles etc or in the case of the `layer0` provider, to support multiple layer0 endpoints.

For example:

```
# aws provider
provider "aws" {
  alias = "east"
  region = "us-east-1"
  # ...
}

# aws provider configured to a west region
provider "aws" {
  alias = "west"
  region = "us-west-1"
  # ...
}
```

This will now allow you to reference aws providers configured to a different region. You can do so by referencing the provider using the naming scheme `TYPE.ALIAS`, which in the above example results in `aws.west`. See [Provider Configuration](https://www.terraform.io/docs/configuration/providers.html) for more information.

```
resource "aws.east_instance" "foo" {
  # ...
}

resource "aws.west_instance" "bar" {
  # ...
}
```

## Part 9: Cleanup

When you're finished with the examples in this guide, run the following destroy command in all the following directories to destroy the Layer0 environment, application and the DynamoDB Table.

Directories:  

 * /example-1  
 * /example-2  
 * /example-3/env-dev  
 * /example-3/env-staging  

`terraform destroy`

!!! Tip "Remote Backend Resources"
	If you created additional resources (S3 bucket and a DynamoDB Table) separately when configuring a [Remote Backend](#part-5-terraform-remote-state), do not forget to delete those if they are no longer needed. You should be able to look at your Terraform configuration file `layer0.tf` to determine the name of the bucket and table.

