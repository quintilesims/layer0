# Layer0 Setup (l0-setup) command-line interface reference

The **l0-setup** application is designed to be used with one of several commands; these commands are detailed in the sections below.

####General Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup** [--version] _command_ [_options_] [_parameters_]</div>
  </div>
</div>

---

##Apply
The **apply** command is used to create and update Layer0 instances.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup apply** [--access\_key=_awsAccessKeyID_] [--secret\_key=_awsSecretAccessKeyID_] [--region=_awsRegion_] [--docker\_token=_dockerToken_] [--vpc=_vpcID_] _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of a Layer0Prefix in an existing AWS stack. To learn more about creating a stack, see the [installation guide](/setup/install.md#part-2-create-an-identity-access-management-user)</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--access\_key=_awsAccessKeyID_</div>
    <div class="divCell">The Access Key ID of an IAM user associated with the AWS stack in which you are creating the Layer0 instance.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--secret\_key=_awsSecretAccessKeyID_</div>
    <div class="divCell">The Secret Access Key ID of an IAM user associated with the AWS stack in which you are creating the Layer0 instance.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--region=_awsRegion_</div>
    <div class="divCell">The AWS region in which the Layer0 instance resides.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--docker\_token=_dockerToken_</div>
    <div class="divCell">A valid d.ims.io Docker token.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--vpc=_vpcID_</div>
    <div class="divCell">The ID of a VPC. If blank, **l0-setup** will create a new VPC.</div>
  </div>
</div>

---

##Backup
The **backup** command is used to back up your Layer0 configuration files to an S3 bucket. This command is most often used when migrating between versions of Layer0. The **backup** command also runs automatically every time you execute the **l0-setup apply** command.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup backup** _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of the Layer0Prefix that you want to back up.</div>
  </div>
</div>

---

##Restore
The **restore** command is used to restore Layer0 configuration files that were previously backed up to an S3 bucket using the [**backup**](#backup) command.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup restore** [--access\_key=_awsAccessKeyID_] [--secret\_key=_awsSecretAccessKeyID_] [--region=_awsRegion_] [--docker\_token=_dockerToken_] _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of the Layer0Prefix that you want to restore.</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--access\_key=_awsAccessKeyID_</div>
    <div class="divCell">The Access Key ID of an IAM user associated with the AWS stack in which the Layer0 instance was created.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--secret\_key=_awsSecretAccessKeyID_</div>
    <div class="divCell">The Secret Access Key ID of an IAM user associated with the AWS stack in which the Layer0 instance was created.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--region=_awsRegion_</div>
    <div class="divCell">The AWS region in which the Layer0 instance resides.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--docker\_token=_dockerToken_</div>
    <div class="divCell">A valid d.ims.io Docker token.</div>
  </div>
</div>

---

##Destroy
The **destroy** command is used to delete a Layer0 configuration.

Before you can run the **destroy** command, you must first delete any existing [loadbalancers](/cli/#loadbalancer-delete), [services](/cli/#service-delete) and [environments](/cli/#environment-delete) that exist in the Layer0 instance, using the **delete** subcommands for each of those components.

!!! warning "Caution"
	Destroying a Layer0 instance cannot be undone; if you created backups of your Layer0 configuration using the <a href="#backup" style="color:#ffffff;">**backup**</a> command, those backups will also be deleted when you run the **destroy** command. The only way to re-deploy a destroyed instance is to use the procedures in the <a href="/setup/install" style="color:#ffffff;">installation guide</a> to rebuild it.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup destroy** _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of a Layer0Prefix in an existing AWS stack.</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--force=false</div>
    <div class="divCell">When specified, destroy confirmation prompts will not be shown.</div>
  </div>
</div>

---

##Endpoint
The **endpoint** command is used to look up the details of a Layer0 endpoint so that you can export them to your shell.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup endpoint** [-iq] [-s _syntax_] _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of a Layer0Prefix for which you want to view endpoint information.</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">-s, --syntax="bash"</div>
    <div class="divCell">Show commands using the specified syntax (bash, powershell, cmd)</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">-i, --insecure=false</div>
    <div class="divCell">Allow incomplete SSL configuration. This option is not recommended for production use.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">-q, --quiet=false</div>
    <div class="divCell">Silence CLI and API version mismatch warning messages</div>
  </div>
</div>

---

##VPC

The **vpc** command is used to look up details from a Virtual Private Cloud (VPC) instance.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup vpc** [--access\_key=_awsAccessKeyID_] [--secret\_key=_awsSecretAccessKeyID_] [--region=_awsRegion_] [--docker\_token=_dockerToken_] _prefixName_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of a Layer0Prefix for which you want to view VPC information.</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--access\_key=_awsAccessKeyID_</div>
    <div class="divCell">The Access Key ID of an IAM user associated with the AWS stack in which you are creating the Layer0 instance.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--secret\_key=_awsSecretAccessKeyID_</div>
    <div class="divCell">The Secret Access Key ID of an IAM user associated with the AWS stack in which you are creating the Layer0 instance.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--region=_awsRegion_</div>
    <div class="divCell">The AWS region in which the Layer0 instance resides.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--docker\_token=_dockerToken_</div>
    <div class="divCell">A valid d.ims.io Docker token.</div>
  </div>
</div>

---

##Terraform
The **terraform** command is used to issue terraform commands directly to a Layer0 instance. Some terraform commands enable functionality that is not a standard part of the **l0-setup** application.

For more information about the capabilities and syntax of terraform commands, see the [Terraform Commands (CLI) Documentation](https://www.terraform.io/docs/commands/index.html).

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0-setup terraform** _prefixName_ _terraformArguments_</div>
  </div>
</div>

###Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_prefixName_</div>
    <div class="divCell">The name of a Layer0Prefix in which you want to execute terraform commands.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_terraformArguments_</div>
    <div class="divCell">The terraform arguments to pass to Layer0. Arguments passed in this way must follow the syntax specified in the [Terraform Commands (CLI) Documentation](https://www.terraform.io/docs/commands/index.html).</div>
  </div>
</div>

---

##Migrate

!!! note
	The **migrate** command has been deprecated and should not be used. This command will be removed from future versions of l0-setup.
