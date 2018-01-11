# Layer0 Setup Reference
The Layer0 Setup application (commonly called `l0-setup`), is used for administrative tasks on Layer0 instances.

##Global options

`l0-setup` can be used with one of several commands: [init](#init), [plan](#plan), [apply](#apply), [list](#list), [push](#push), [pull](#pull), [endpoint](#endpoint), [destroy](#destroy), [upgrade](#upgrade), and [set](#set). These commands are detailed in teh sections below. There are, however, some global paramters that you may specify whenever using `l0-setup`

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">`l0-setup [global options] command params`</div>
  </div>
</div>

###Global options
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">`-l value, --log value`</div>
    <div class="divCell">The log level to display on the console when you run commands. (default: info)</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--version`</div>
    <div class="divCell">Display the version number of the `l0-setup` application.</div>
  </div>
</div>

---

## Init
The `init` command is used to initialize or reconfigure a Layer0 instance. 
This command will prompt the user for inputs required to create/update a Layer0 instance. 
Each of the inputs can be specified through an optional flag.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">`l0-setup init [--docker-path path | --module-source path | --version version | --aws-access-key access_key | --aws-secret-key secret_key | --aws-region region] instanceName`</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">`--docker-path`</div>
    <div class="divCell">Path to docker config.json file. This is used to include private Docker Registry authentication for this Layer0 instance.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--module-source`</div>
    <div class="divCell">The source input variable is the path to the Terraform Layer0. By default, this points to the Layer0 github repository. Using values other than the default may result in undesired consequences.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--version`</div>
    <div class="divCell">The version input variable specifies the tag to use for the Layer0 Docker images: `quintilesims/l0-api` and `quintilesims/l0-runner`.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--aws-access-key`</div>
    <div class="divCell">The access_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--aws-secret-key`</div>
    <div class="divCell">The secret_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--aws-ssh-key-pair`</div>
    <div class="divCell">The ssh_key_pair input variable specifies the name of the ssh key pair to include in EC2 instances provisioned by Layer0. This key pair must already exist in the AWS account.  The names of existing key pairs can be found in the EC2 dashboard.</div>
  </div>
</div>

---

## Plan
The `plan` command is used to show the planned operation(s) to run during the next `apply` on a Layer0 instance without actually executing any actions

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">`l0-setup plan instanceName`</div>
  </div>
</div>

---

## Apply
The `apply` command is used to create and update Layer0 instances. Note that the default behavior of apply is to push the layer0 configuration to an S3 bucket unless the `--push=false` flag is set to false. Pushing the configuration to an S3 bucket requires aws credentials which if not set via the optional `--aws-*` flags, are read from the environment variables or a credentials file. 

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">`l0-setup apply [--quick | --push=false | --aws-access-key | --aws-secret-key | --aws-region] instanceName`</div>
  </div>
</div>

###Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">`--quick`</div>
    <div class="divCell">Skips verification checks that normally run after `terraform apply` has completed</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--push=false`</div>
    <div class="divCell">Skips uploading local Layer0 configuration files to an S3 bucket</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--aws-access-key`</div>
    <div class="divCell">The Access Key ID portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">`--aws-secret-key`</div>
    <div class="divCell">The Secret Access Key portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI.</div>
  </div>
</div>

---

## List
The `list` command is used to list local and remote Layer0 instances.

### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">`l0-setup list [--local | --remote | --aws-access-key | --aws-secret-key]`</div>
  </div>
</div>

### Options
* `-l, --local` - Show local Layer0 instances. This value is true by default.
* `-r, --remote` - Show remote Layer0 instances. This value is true by default. 
* `--aws-access-key` - The Access Key ID portion of an AWS Access Key that has permissions to list S3 buckets. 
If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-secret-key` - The Secret Access Key portion of an AWS Access Key that has permissions to list S3 buckets. 
If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 

---
## Push
The **push** command is used to back up your Layer0 configuration files to an S3 bucket.

### Usage
```
$ l0-setup push [options] <instance_name> 
```

### Options
* `--aws-access-key` - The Access Key ID portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-secret-key` - The Secret Access Key portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-region` - The region of the Layer0 instance. The default value is `us-west-2`. 

---
## Pull
The **pull** command is used copy Layer0 configuration files from an S3 bucket.

### Usage
```
$ l0-setup pull [options] <instance_name> 
```

### Options
* `--aws-access-key` - The Access Key ID portion of an AWS Access Key that has permissions to pull to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-secret-key` - The Secret Access Key portion of an AWS Access Key that has permissions to pull to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-region` - The region of the Layer0 instance. The default value is `us-west-2`. 

---
## Endpoint
The **endpoint** command is used to show environment variables used to connect to a Layer0 instance

### Usage
```
$ l0-setup endpoint [options] <instance_name> 
```

### Options
* `-i, --insecure` - Show environment variables that allow for insecure settings
* `-d, --dev` - Show environment variables that are required for local development
* `-s --syntax` - Choose the syntax to display environment variables 
(choices: `bash`, `cmd`, `powershell`) (default: `bash`)

---
## Destroy
The **destroy** command is used to destroy all resources associated with a Layer0 instance.

!!! warning "Caution"
	Destroying a Layer0 instance cannot be undone; if you created backups of your Layer0 configuration using the **push** command, those backups will also be deleted when you run the **destroy** command.

### Usage
```
$ l0-setup destroy [options] <instance_name> 
```

### Options
* `--force` - Skips confirmation prompt


---
## Upgrade
The **upgrade** command is used to upgrade a Layer0 instance to a new version.
You will need to run an **apply** after this command has completed. 


### Usage
```
$ l0-setup upgrade [options] <instance_name> <version>
```

### Options
* `--force` - Skips confirmation prompt


---
## Set
The **set** command is used set input variable(s) for a Layer0 instance's Terraform module.
This command can be used to shorthand the **init** and **upgrade** commands, 
and can also be used with custom Layer0 modules. 
You will need to run an **apply** after this command has completed. 

### Usage
```
$ l0-setup set [options] <instance_name>
```

**Example Usage**
```
$ l0-setup set --input username=admin --input password=pass123 mylayer0
```

### Options
* `--input` - Specify an input using `key=val` format
