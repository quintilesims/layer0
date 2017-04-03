# Deployment guide: Guestbook sample application
In this example, you will learn how you can use Terraform to create a Layer0 service as well as supporting infrastructure in the form of a persistent data store. This example does not cover the details of application being deployed but focuses on how you can combine Layer0 and other Terraform providers as part of a Terraform workflow.

## Before you start
In order to complete the procedures in this section, you must have following installed and configured correctly:

 * Layer0
 * Terraform
 * Layer0 Terraform Provider

Layer0 v0.8.4 or later installed and configured. If you have not already configured Layer0, see the [Layer0 installation guide](/setup/install). If you are running an older version of Layer0, see the [Layer0 upgrade instructions](/setup/upgrade#upgrading-older-versions-of-layer0).

See the [Terraform installation guide](/reference/terraform-plugin#install) to install Terraform and the Layer0 Terraform Plugin.

---

## Deploy with Terraform
Using Terraform, you will deploy a simple guestbook application, which is backed by AWS DynamoDb Table, for persistant storage. The terraform configuration file will use both the Layer0 and AWS Terraform providers, to deploy the guestbook application and provision a new DynamoDb Table.

## Part 1: Clone the Layer0 examples repository
Run this command to clone the `quintilesims/layer0-examples` repository

<ul>
  <li class="command">git clone https://github.com/quintilesims/layer0-examples.git</li>
</ul>

Inside a terminal window, navigate to the `terraform-beyond-layer0 folder`. You should find the following files once you are in the folder. We use these files to set up a Layer0 envrionment and deploy AWS resources with Terraform:

|Filename|Purpose|
|----|----|
|[terraform.tfvars](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/terraform.tfvars)|Variables specific to the environment and guestbook application|
|[Dockerrun.aws.json](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/Dockerrun.aws.json)|Template for running the guestbook application in a Layer0 environment|
|[layer0.tf](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/layer0.tf)|Provision Layer0 resources|

## Part 2: Terraform Plan
Before deploying, we can run the follwing command to see what changes Terraform will make to your infrastructure should you go ahead and apply. If you had any errors in your layer0.tf file, running `terraform plan` would output those errors so that you can address them.

<ul>
  <li class="command">terraform plan</li>
</ul>

```
+ aws_dynamodb_table.guestbook
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
Run the follwing command to begin the deploy process. Terraform will prompt you for configuration values that it does not have.

<ul>
  <li class="command">terraform apply</li>
</ul>

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
Terraform using the AWS provider, provisions a new DynamoDb Table (with the name you have configured) and configures environment variables for the layer0 application to use the newly provisioned table, and deploys the application into a Layer0 environment using the Layer0 provider.

Looking at excerpt of [layer0.tf](https://github.com/quintilesims/layer0-examples/blob/master/guestbook-db/layer0.tf) file, we can see the following definitions:

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

Note the resources definitions for `aws_dynamodb_table` and `layer0_deploy`. To configure the guestbook application to use the provisioned dynamodb table, we reference the `name` property from the dynamodb definition `table_name = "${aws_dynamodb_table.guestbook.name}"`. 

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

The Layer0 configuration referencing the Aws DynamoDb configuration `table_name = "${aws_dynamodb_table.guestbook.name}"`, infers an implicit dependency. Before Terraform creates the infrastructre, it will use this information to order the resources created and also create resources in parallel where there are no dependencies. In this example, the AWS DynamoDb table will be created before the Layer0 deploy. See [Terraform Dependencies](https://www.terraform.io/intro/getting-started/dependencies.html) for more information.

## Part 4: Terraform Destroy
When you're finished with the example run the following command in the same directory to destroy the Layer0 environment, application and the DynamoDb Table.

<ul>
  <li class="command">terraform destroy</li>
</ul>
