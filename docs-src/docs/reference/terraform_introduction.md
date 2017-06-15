# Introduction to Terraform


## What does Terraform do?

Terraform is a powerful orchestration tool for creating, updating, deleting, and otherwise managing infrastructure in an easy-to-understand, declarative manner.
Terraform's [documentation](https://www.terraform.io/intro/index.html) is very good, but at a glance:

**Be Declarative -**
Specify desired infrastructure results in Terraform (`*.tf`) files, and let Terraform do the heavy work of figuring out how to make that specification a reality.

**Scry the Future -**
Use `terraform plan` to see a list of everything that Terraform _would_ do without actually making those changes.

**Version Infrastructure -**
Check Terraform files into a VCS to track changes to and manage versions of your infrastructure.


## Why Terraform?

Why did we latch onto Terraform instead of something like CloudFormation?

**Cloud-Agnostic -**
Unlike CloudFormation, Terraform is able to incorporate different [resource providers](https://www.terraform.io/docs/providers/index.html) to manage infrastructure across multiple cloud services (not just AWS).

**Custom Providers -**
Terraform can be extended to manage tools that don't come natively through use of custom providers.
We wrote a [Layer0 provider](/reference/terraform-plugin) so that Terraform can manage Layer0 resources in addition to tools and resources and infrastructure beyond Layer0's scope.

Terraform has some [things to say](https://www.terraform.io/intro/vs/index.html) on the matter as well.


## Advantages Versus Layer0 CLI?

Why should you move from using (or scripting) the Layer0 CLI directly?

**Reduce Fat-Fingering Mistakes -**
Creating Terraform files (and using `terraform plan`) allows you to review your deployment and catch errors.
Executing Layer0 CLI commands one-by-one is tiresome, non-transportable, and a process ripe for typos.

**Go Beyond Layer0 -**
Retain the benefits of leveraging Layer0's concepts and resources using our [provider](/reference/terraform-plugin), but also gain the ability to orchestrate resources and tools beyond the CLI's scope.


## How do I get Terraform?

Check out Terraform's [documentation](https://www.terraform.io/intro/getting-started/install.html) on the subject.
