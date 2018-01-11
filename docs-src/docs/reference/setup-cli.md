# Layer0 Setup Reference
The Layer0 Setup application (commonly called `l0-setup`), is used for administrative tasks on Layer0 instances.

## Global options

`l0-setup` can be used with one of several commands: [init](#init), [plan](#plan), [apply](#apply), [list](#list), [push](#push), [pull](#pull), [endpoint](#endpoint), [destroy](#destroy), [upgrade](#upgrade), and [set](#set). These commands are detailed in teh sections below. There are, however, some global paramters that you may specify whenever using `l0-setup`

### Usage
```
l0-setup [global options] command [command options] params
```

### Global options
* `-l value, --log value` - The log level to display on the console when you run commands. (default: info)
* `--version` - Display the version number of the `l0-setup` application.

---

## Init
The `init` command is used to initialize or reconfigure a Layer0 instance. 
This command will prompt the user for inputs required to create/update a Layer0 instance. 
Each of the inputs can be specified through an optional flag.

### Usage
```
l0-setup init [--docker-path path | --module-source path | 
    --version version | --aws-region region | --aws-access-key accessKey | 
    --aws-secret-key secretKey] instanceName`
```

### Optional arguments
* `--docker-path` - Path to docker config.json file. This is used to include private Docker Registry authentication for this Layer0 instance.
* `--module-source` - The source input variable is the path to the Terraform Layer0. By default, this points to the Layer0 github repository. Using values other than the default may result in undesired consequences.
* `--version` - The version input variable specifies the tag to use for the Layer0 Docker images: `quintilesims/l0-api` and `quintilesims/l0-runner`.
* `--aws-ssh-key-pair` - The ssh_key_pair input variable specifies the name of the ssh key pair to include in EC2 instances provisioned by Layer0. This key pair must already exist in the AWS account.  The names of existing key pairs can be found in the EC2 dashboard.
* `--aws-access-key` - The access_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.
* `--aws-secret-key` - The secret_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.

---

## Plan
The `plan` command is used to show the planned operation(s) to run during the next `apply` on a Layer0 instance without actually executing any actions

### Usage
```
l0-setup plan instanceName
```

---

## Apply
The `apply` command is used to create and update Layer0 instances. Note that the default behavior of apply is to push the layer0 configuration to an S3 bucket unless the `--push=false` flag is set to false. Pushing the configuration to an S3 bucket requires aws credentials which if not set via the optional `--aws-*` flags, are read from the environment variables or a credentials file. 

### Usage
```
l0-setup apply [--quick | --push=false | --aws-access-key accessKey | 
    --aws-secret-key secretKey] instanceName`
```

### Optional arguments
* `--quick` - Skips verification checks that normally run after `terraform apply` has completed
* `--push=false` - Skips uploading local Layer0 configuration files to an S3 bucket
* `--aws-access-key` - The access_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.
* `--aws-secret-key` - The secret_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.
---

## List
The `list` command is used to list local and remote Layer0 instances.

### Usage
```
l0-setup list [--local=false | --remote=false | --aws-access-key accessKey | 
    --aws-secret-key secretKey]
```

### Optional arguments
* `-l, --local` - Show local Layer0 instances. This value is true by default.
* `-r, --remote` - Show remote Layer0 instances. This value is true by default. 

---
## Push
The `push` command is used to back up your Layer0 configuration files to an S3 bucket.

### Usage
```
l0-setup push [--aws-access-key accessKey | 
    --aws-secret-key secretKey] instanceName
```

### Optional arguments
* `--aws-access-key` - The access_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.
* `--aws-secret-key` - The secret_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.

---
## Pull
The `pull` command is used copy Layer0 configuration files from an S3 bucket.

### Usage
```
l0-setup pull [--aws-access-key accessKey | 
    --aws-secret-key secretKey] instanceName
```

### Optional arguments
* `--aws-access-key` - The access_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.
* `--aws-secret-key` - The secret_key input variable is used to provision the AWS resources required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key. It is recommended this key has the `AdministratorAccess` policy.

---
## Endpoint
The `endpoint` command is used to show environment variables used to connect to a Layer0 instance

### Usage
```
l0-setup endpoint [-i | -d | -s syntax] instanceName
```

### Optional arguments
* `-i, --insecure` - Show environment variables that allow for insecure settings
* `-d, --dev` - Show environment variables that are required for local development
* `-s --syntax` - Choose the syntax to display environment variables 
(choices: `bash`, `cmd`, `powershell`) (default: `bash`)

---
## Destroy
The `destroy` command is used to destroy all resources associated with a Layer0 instance.

!!! danger "Caution"
	Destroying a Layer0 instance cannot be undone. If you created backups of your Layer0 configuration using the `push` command, those backups will also be deleted when you run the `destroy` command.

### Usage
```
l0-setup destroy [--force] instanceName
```

### Optional arguments
* `--force` - Skips confirmation prompt

---
## Upgrade
The `upgrade` command is used to upgrade a Layer0 instance to a new version.
You will need to run an `apply` after this command has completed. 


### Usage
```
l0-setup upgrade [--force] instanceName version
```

### Optional arguments
* `--force` - Skips confirmation prompt

---
## Set
The `set` command is used set input variable(s) for a Layer0 instance's Terraform module.
This command can be used to shorthand the `init` and `upgrade` commands, and can also be used with custom Layer0 modules. 
You will need to run an `apply` after this command has completed. 

### Usage
```
l0-setup set [--input key=value] instanceName
```

### Options
* `--input key=val` - Specify an input using `key=val` format

### Example Usage
```
l0-setup set --input username=admin --input password=pass123 mylayer0
```

