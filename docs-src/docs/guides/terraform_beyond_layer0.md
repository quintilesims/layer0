# Deployment guide: Terraform beyond Layer0
 In this example, we'll learn how you can use Terraform to create a Layer0 service as well as a persistent data store. The main goal of this example is to explore how you can combine Layer0 with other Terraform providers.

## Before you start
In order to complete the procedures in this section, you must have the following installed and configured correctly:

 * Layer0 v0.8.4 or later
 * Terraform v0.9.0 or later
 * Layer0 Terraform Provider

If you have not already configured Layer0, see the [Layer0 installation guide](/setup/install). If you are running an older version of Layer0, see the [Layer0 upgrade instructions](/setup/upgrade#upgrading-older-versions-of-layer0).

See the [Terraform installation guide](/reference/terraform-plugin#install) to install Terraform and the Layer0 Terraform Plugin.

---

## Deploy with Terraform
Using Terraform, you will deploy a simple guestbook application, which is backed by AWS DynamoDB Table, for persistant storage. The terraform configuration file will use both the Layer0 and AWS Terraform providers, to deploy the guestbook application and provision a new DynamoDB Table.

## Part 1: Clone the Layer0 examples repository
Run this command to clone the `quintilesims/layer0-examples` repository

`git clone https://github.com/quintilesims/layer0-examples.git`  

Inside a terminal window, navigate to the `terraform-beyond-layer0 folder`. You should find the following files once you are in the folder. We use these files to set up a Layer0 environment and deploy AWS resources with Terraform:

|Filename|Purpose|
|----|----|
|[terraform.tfvars](https://github.com/quintilesims/layer0-examples/blob/master/terraform-beyond-layer0/terraform.tfvars)|Variables specific to the environment and guestbook application|
|[Dockerrun.aws.json](https://github.com/quintilesims/layer0-examples/blob/master/terraform-beyond-layer0/Dockerrun.aws.json)|Template for running the guestbook application in a Layer0 environment|
|[layer0.tf](https://github.com/quintilesims/layer0-examples/blob/master/terraform-beyond-layer0/layer0.tf)|Provision Layer0 resources and AWS resources|

## Part 2: Terraform Plan
Before deploying, we can run the follwing command to see what changes Terraform will make to your infrastructure should you go ahead and apply. If you had any errors in your layer0.tf file, running `terraform plan` would output those errors so that you can address them. Also, Terraform will prompt you for configuration values that it does not have.

`terraform plan`

!!! Note
	There are a few ways to configure Terraform so that you don't have to keep entering these values every time you run a Terraform command (editing the `terraform.tfvars` file, or exporting evironment variables like `TF_VAR_endpoint` and `TF_VAR_token`, for example). See the [Terraform Docs](https://www.terraform.io/docs/configuration/variables.html) for more.

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
  It may take a few minutes for the guestbook service to launch and the load balancer to become available. During that time you may get HTTP 503 errors when making HTTP requests against the load balancer URL.

Terraform will set up the entire environment for you and then output a link to the application's load balancer.

### What's happening
Terraform using the AWS provider, provisions a new DynamoDB Table (with the name you have configured) and configures environment variables for the Layer0 application to use the newly provisioned table, and deploys the application into a Layer0 environment using the Layer0 provider.

Looking at an excerpt of [layer0.tf](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/layer0.tf) file, we can see the following definitions:

```
resource "aws_DynamoDB_table" "guestbook" {
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
    table_name = "${aws_DynamoDB_table.guestbook.name}"
  }
}
```

Note the resource definitions for `aws_dynamodb_table` and `layer0_deploy`. To configure the guestbook application to use the provisioned DynamoDB table, we reference the `name` property from the DynamoDB definition `table_name = "${aws_dynamodb_table.guestbook.name}"`. 

This is then used to populate the template fields in our [Dockerrun.aws.json](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/Dockerrun.aws.json) file. 

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

The Layer0 configuration referencing the AWS DynamoDB configuration `table_name = "${aws_DynamoDB_table.guestbook.name}"`, infers an implicit dependency. Before Terraform creates the infrastructre, it will use this information to order the resources created and also create resources in parallel where there are no dependencies. In this example, the AWS DynamoDB table will be created before the Layer0 deploy. See [Terraform Dependencies](https://www.terraform.io/intro/getting-started/dependencies.html) for more information.

## Part 4: Scaling a Layer0 service
The workflow to make changes to your infrastructure generally involves updating your Terraform configuration file followed by a `terraform plan` and `terraform apply`.

### Update the Terraform configuration
Open the file `layer0.tf` in a text editor and make the change to add a `scale` property with a value of `3` to the `layer0_service` section. For more information about the `scale` property, see [Layer0 Terraform Plugin](http://layer0.ims.io/reference/terraform-plugin/#service) documenation. The end result should look like the below:


layer0.tf

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
Once you have updated the following command to understand the changes that you will be making. Note that if you did not specify `scale` it defaults to '1'.

`terraform plan`

Outputs:

```
...

~ layer0_service.guestbook
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

guestbook_url = <guestbook_service_url>
```

To confirm your service has scaled, you can run the following layer0 command. Note desired scale for the guestbook service should be eventually be 3/3.

`l0 service get demo-env:guestbook`

Outputs:

```
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
api           api           api          api           api:3        1/1
guestboebca1  guestbook     demo-env     guestbook     guestbook:6  3/3
```



## Part 5: Terraform Destroy
When you're finished with the example run the following command in the same directory to destroy the Layer0 environment, application and the DynamoDB Table.

`terraform destroy`

---

# Best Practices with Terraform + Layer0

## Part 6: Terraform Remote State

Terraform stores the state of the deployed infrastructure in a local file named `terraform.tfstate` by default. This includes not only information about the resources deployed but also metadata such as resource dependencies. To find out more about why Terraform needs to store state see [Purpose of Terraform State](https://www.terraform.io/docs/state/purpose.html). 

How state is loaded and used for operations such as `terraform apply` is determined by a [Backend](https://www.terraform.io/docs/backends). As mentioned, by default the state is stored locally which is enabled by a "local" backend.

### Remote State
In scenarios where you are working as part of a team to provision and manage a services deployed by Terraform, all the members of the team will need access to the state file to apply new changes. You would also want resiliency against losing the state file. You can cater for provisioning and managing Terraform in a team environment and providing redundancy for the state file by using a remote backend. A remote backend can also provide locking mehcanisms to ensure multiple users can't change resources at the same time. Backends such as [Consul](https://www.terraform.io/docs/backends/types/consul.html) & [S3](https://www.terraform.io/docs/backends/types/s3.html) among others, are backends that you can use to store master Terraform state in a remote location.

To configure a remote backend, append the `terraform` section below to your terraform file `layer0.tf`. Update the `path` property to a GUID to avoid a potential conflict as demo.consul.io is a public consul endpoint.

```
terraform {
  backend "s3" {
    bucket     = "<my-bucket-name>"
    key        = "demo-env/remote-backend/terraform.tfstate"
    region     = "us-west-2"
  }
}
```

Once you have modified `layer0.tf`, you will need to initailize the newly configured backend by running the following command.

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
As you are configuring a backend for the first time, Terraform will give you an option to migrate your state to the new backend. From now on, any further changes to your infrastructure made by Terraform will result in the remote state file being updated. For more information about see (Terraform backends)[https://www.terraform.io/docs/backends/index.html].

A new team member can use the `layer0.tf` from their own machine without obtaining a copy of the state file `terraform.tfstate` as the configuration points to a remote backend where the state file will be retrieved from.

### Locking
Not all remote backends support locking (locking ensures only one person is able to change the tfstate at a time). The `S3` backend we used earlier in the example also supports locking which is disabled by default. To enable locking, you need to specify `locking_table` property with the name of an existing DyanmoDB table. The DynamoDB table also needs primary key named `LockID` of type `String`.


### Security
A Terraform state file is written in plain text. This can lead to a situation where deploying resources that require sensitive data can result in the sensitive data being stored in the state file. To minimize exposure of senstive data, you can enable [server side encryption](https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html) of the state file by adding property `encrypt` set to `true`.

This will ensure that the file is encrypted in S3 and by using a remote backend, you will also have the added benefit of the state file not being persisted to disk locally as it will only ever be held in memory by Terraform.

For securing the state file further, you can also enable access logging on the S3 bucket you are using for the remote backend which can help track down invalid access should it occur.

## Part 8: Terraform Configuration structure

There are many different approaches to setup your Terraform file structure. There is no single best precribed way to structure your files. Whatever approach you take needs to be catered for your needs. Keeping that in mind; the file structure for [Terraform beyond Layer0 example](https://github.com/quintilesims/layer0-examples/blob/master/terraform-beyond-layer0) is:

* root  
  + main.tf  
  + variables.tf  
  + output.tf  
    + modules  
      + guestbook_service  
        + main.tf  
        + variables.tf  
        + output.tf
      + service2
      + service3

Here we are making use of Terraform [Modules](https://www.terraform.io/docs/modules/index.html). Modules in Terraform are self-contained packages of Terraform configurations that are managed as a group. Modules are used to create reusable components in Terraform as well as for basic code organization. In this example, we are using modules to separate each service and making it consumable as a module.

If you wanted to add a new service, you can create new folder service folder inside the ./modules. If you wanted to you could even run multiple copies of the same service. See here for more information about [Creating Modules](https://www.terraform.io/docs/modules/create.html).

When creating a module, ensure that resources you are creating are prefixed with the environment and the module's name variable to ensure your resources are unique for each layer0 environment and each reference to a module.
```
resource "layer0_load_balancer" "guestbook" {
  name        = "${var.name}_guestbook_lb"
  environment = "${var.layer0_environment_id}"

  port {
    host_port      = 80
    container_port = 80
    protocol       = "http"
  }
}

resource "aws_dynamodb_table" "guestbook" {
  name           = "${var.layer0_environment_name}_${var.name}_${var.table_name}"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}
```

Also see the below repositories for ideas on different ways you can organize your project:
* [Terraform Community Modules](https://github.com/terraform-community-modules)
* [Best Pratices Ops](https://github.com/hashicorp/best-practices)
