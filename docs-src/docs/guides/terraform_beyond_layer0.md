# Deployment guide: Terraform beyond Layer0
In this example, you will learn how you can use Terraform to create a Layer0 service as well as supporting infrastructure in the form of a persistent data store. This example does not cover the details of the application being deployed but focuses on how you can combine Layer0 and other Terraform providers as part of a Terraform workflow.

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

## Part 4: Terraform Remote State

Terraform stores the state of the deployed infrastructure in a local file named `terraform.tfstate` by default. This includes not only information about the resources deployed but also metadata such as resource dependencies. To find out more about why Terraform needs to store state see [Purpose of Terraform State](https://www.terraform.io/docs/state/purpose.html). 

How state is loaded and used for operations such as `terraform apply` is determined by a [Backend](https://www.terraform.io/docs/backends). As mentioned, by default the state is stored locally which is enabled by a "local" backend.

### Remote State
In scenarios where you are working as part of a team to provision and manage a services deployed by Terraform, all the members of the team will need access to the state file to apply new changes. You would also want resiliency against losing the state file. You can cater for provisioning and managing Terraform in a team environment and providing redundancy for the state file by using a remote backend. Backends such as [Consul](https://www.terraform.io/docs/backends/types/consul.html) & [S3](https://www.terraform.io/docs/backends/types/s3.html) among others, are backends that you can use to store master Terraform state in a remote location.

To configure a remote backend, append the `terraform` section below to your terraform file `layer0.tf`. Update the `path` property to a GUID to avoid a potential conflict as demo.consul.io is a public consul endpoint.

```
terraform {
  backend "consul" {
    address = "demo.consul.io"
    path    = "a-unique-random-string"
    lock    = false
  }
}
```

Once you have modified `layer0.tf`, you will need to initailize the newly configured backend by running the following command.

`terraform init`

Outputs:

```
Initializing the backend...

Do you want to copy state from "local" to "consul"?
  Pre-existing state was found in "local" while migrating to "consul". An existing
  non-empty state exists in "consul". The two states have been saved to temporary
  files that will be removed after responding to this query.
  
  One ("local"): /var/folders/n0/19h9crxn2v75g70bl0txdj1wm05615/T/terraform529524247/1-local.tfstate
  Two ("consul"): /var/folders/n0/19h9crxn2v75g70bl0txdj1wm05615/T/terraform529524247/2-consul.tfstate
  
  Do you want to copy the state from "local" to "consul"? Enter "yes" to copy
  and "no" to start with the existing state in "consul".

  Enter a value: 

```

Go ahead and enter: `yes`.

```

Successfully configured the backend "consul"! Terraform will automatically
use this backend unless the backend configuration changes.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your environment. If you forget, other
commands will detect it and remind you to do so if necessary.

```

### What's happening
As you are configuring a backend for the first time, Terraform will give you an option to migrate your state to the new backend. From now on, any further changes to your infrastructure made by Terraform will result in the remote state file being updated. For more information about see (Terraform backends)[https://www.terraform.io/docs/backends/index.html].

A new team member can use the `layer0.tf` from their own machine without obtaining a copy of the state file `terraform.tfstate` as the configuration points to a remote backend where the state file will be retrieved from.

The command `terraform init` must be called:
 * on any new environment that configures a backend
 * on any change of the backend configuration (including type of backend)
 * on removing backend configuration completely

## Part 5: Scaling a Layer0 service
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
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

layer0_environment.demo: Refreshing state... (ID: demoenvbb9f6)
data.template_file.guestbook: Refreshing state...
layer0_deploy.guestbook: Refreshing state... (ID: guestbook.6)
layer0_load_balancer.guestbook: Refreshing state... (ID: guestbo43ab0)
layer0_service.guestbook: Refreshing state... (ID: guestboebca1)
The Terraform execution plan has been generated and is shown below.
Resources are shown in alphabetical order for quick scanning. Green resources
will be created (or destroyed and then created if an existing resource
exists), yellow resources are being changed in-place, and red resources
will be destroyed. Cyan entries are data sources to be read.

Note: You didn't specify an "-out" parameter to save this plan, so when
"apply" is called, Terraform can't guarantee this is what will execute.

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

`l0 service get \*`

Outputs:

```
SERVICE ID    SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYMENTS  SCALE
api           api           api          api           api:3        1/1
guestboebca1  guestbook     demo-env     guestbook     guestbook:6  3/3
```



## Part 6: Terraform Destroy
When you're finished with the example run the following command in the same directory to destroy the Layer0 environment, application and the DynamoDB Table.

`terraform destroy`

## Best Practices with Terraform + Layer0


