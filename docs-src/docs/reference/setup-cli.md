# Layer0 Setup Reference
The Layer0 Setup application (commonly called **l0-setup**), is used to provision, update, and destroy Layer0 instances.

---
## General Usage
You can use the `-h, --help` command to get generate information about the `l0-setup` tool:

---
## Init
The **init** command is used to initialize or reconfigure a Layer0 instance. 
This command will prompt the user for inputs required to create/update a Layer0 instance. 
Each of the inputs can be specified through an optional flag.

### Usage
```
$ l0-setup init [options] <instance_name> 
```

### Options
* `--docker-path` - Path to docker config.json file. 
This is used to include private Docker Registry authentication for this Layer0 instance.
* `--module-source` - The source input variable is the path to the Terraform Layer0. 
By default, this points to the Layer0 github repository. 
Using values other than the default may result in undesired consequences.
* `--version` - The version input variable specifies the tag to use for the Layer0
Docker images: `quintilesims/l0-api`.
* `--aws-access-key` - The access_key input variable is used to provision the AWS resources
required for Layer0. 
This corresponds to the Access Key ID portion of an AWS Access Key.
It is recommended this key has the `AdministratorAccess` policy. 
* `--aws-secret-key` The secret_key input variable is used to provision the AWS resources
required for Layer0. 
This corresponds to the Secret Access Key portion of an AWS Access Key.
It is recommended this key has the `AdministratorAccess` policy.
* `--aws-region` - The region input variable specifies which region to provision the
AWS resources required for Layer0. The following regions can be used:
    - us-west-1
    - us-west-2
    - us-east-1
    - eu-west-1


* `--aws-ssh-key-pair` - The ssh_key_pair input variable specifies the name of the
ssh key pair to include in EC2 instances provisioned by Layer0. 
This key pair must already exist in the AWS account. 
The names of existing key pairs can be found in the EC2 dashboard.

---
## Plan
The **plan** command is used to show the planned operation(s) to run during the next `apply` on a Layer0 instance without actually executing any actions

### Usage
```
$ l0-setup plan <instance_name> 
```

### Options
There are no options for this command

---
## Apply
The **apply** command is used to create and update Layer0 instances. Note that the default behavior of apply is to push the layer0 configuration to an S3 bucket unless the `--push=false` flag is set to false. Pushing the configuration to an S3 bucket requires aws credentials which if not set via the optional `--aws-*` flags, are read from the environment variables or a credentials file. 

### Usage
```
$ l0-setup apply [options] <instance_name> 
```

### Options
* `--quick` - Skips verification checks that normally run after `terraform apply` has completed
* `--push` - Skips uploading local Layer0 configuration files to an S3 bucket
* `--aws-access-key` - The Access Key ID portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-secret-key` - The Secret Access Key portion of an AWS Access Key that has permissions to push to the Layer0 instances's S3 bucket. If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-region` - The region of the Layer0 instance. The default value is `us-west-2`. 

---
## List
The **list** command is used to list local and remote Layer0 instances.

### Usage
```
$ l0-setup list [options]
```

### Options
* `-l, --local` - Show local Layer0 instances. This value is true by default.
* `-r, --remote` - Show remote Layer0 instances. This value is true by default. 
* `--aws-access-key` - The Access Key ID portion of an AWS Access Key that has permissions to list S3 buckets. 
If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-secret-key` - The Secret Access Key portion of an AWS Access Key that has permissions to list S3 buckets. 
If not specified, the application will attempt to use any AWS credentials used by the AWS CLI. 
* `--aws-region` - The region to list S3 buckets. The default value is `us-west-2`. 

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
