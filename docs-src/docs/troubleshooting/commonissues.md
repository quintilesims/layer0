#Common issues and their solutions

##Manually deleting a Layer0 instance from AWS

Sometimes, your Layer0 instance might get into an unresponsive and unrecoverable state. This section describes the AWS resources that are created when you create a Layer0 instance and what resources needs to be removed by hand in such an event.

!!! Tip
    An unresponsive state means that the Layer0 API isn't responding to typical commands (from the CLI or through Terraform). If Layer0 is unresponsive, it might still be recoverable. You should try these tips first if you're running into problems.

    [l0-setup init](../reference/setup-cli#init) and [l0-setup apply](../reference/setup-cli#apply)

    Terraform might be able to refresh the resources it created when the Layer0 instance was first created. You can use `l0-setup init` and `l0-setup apply` to double check your instance's initial settings and attempt to rebuild core VPC or API resources if they were accidentally altered.

    Use [l0-setup endpoint](../reference/setup-cli#endpoint) to double check your instance's environment variables.

    Are your `LAYER0_*` environment variables set correctly? Is the URL of your Layer0 API correct? Are the AWS Access Key and Secret Access Key values still valid?

    If you are unable to recover your Layer0 instance's state and need to delete your instance and start over, continue reading this guide.

Each instance of `<prefix>` is the name of the Layer0 instance you specified when using `l0-setup`

* VPC
    * Name: `l0-<prefix>`
    * Subnets
        * 3 Public subnets: `l0-<prefix>-subnet-public-<region & availability zone>`
        * 3 Private subnets: `l0-<prefix>-subnet-private-<region & availability zone>`
    * Route Tables
        * A blank default
        * `l0-<prefix>-rt-private`
        * `l0-<prefix>-rt-public`
    * Internet Gateway: `l0-<prefix>-igw`
    * NAT Gateway: nameless but associated with the VPC
    * Network ACL: nameless but associated with the VPC
    * Security Groups
        * `default`
        * `l0-<prefix>-api-lb`
        * `l0-<prefix>-api-env`

* EC2
    * Auto Scaling Group: `l0-<prefix>-api`
    * Launch Configuration: `l0-<prefix>-api-<timestamp>`
    * Instances
        * 2 EC2 instances named `l0-<prefix>-api`
    * Load Balancer: `l0-<prefix>-api`

* EC2 Container Service
    * Cluster: `l0-<prefix>-api`
    * Task Definition: `l0-<prefix>-api`

* IAM
    * Group: `l0-<prefix>`
    * User: `l0-<prefix>-user`
    * Role: `l0-<prefix>-ecs-role`
    * Instance profile: `l0-<prefix>-ecs-instance-profile`
    * Server certificate: `l0-<prefix>-api`

* S3
    * Bucket: `layer0-<prefix>-<accountnumber>`

* CloudWatch
    * Log Group: `l0-<prefix>`

* DynamoDB
    * Tables
        * `l0-<prefix>-lock`
        * `l0-<prefix>-tags`


Most resources can be removed through the AWS Console, but some need to be removed from the AWS CLI.

**IAM Instance Profile**

`aws iam list-instance-profiles` to list

`aws iam delete-instance-profile --instance-profile-name [name]` to delete

**IAM Server Certificate**

`aws iam list-server-certificates` to list

`aws iam delete-server-certificate --server-certificate-name [name]`


##"Connection refused" error when executing Layer0 commands

When executing commands using the Layer0 CLI, you may see the following error message: 

`Get http://localhost:9090/command/: dial tcp 127.0.0.1:9090: connection refused`

Where `command` is the Layer0 command you are trying to execute.

This error indicates that your Layer0 environment variables have not been set for the current session. See the ["Connect to a Layer0 Instance" section](../setup/install/#part-4-connect-to-a-layer0-instance) of the Layer0 installation guide for instructions for setting up your environment variables.

---

## "Invalid Dockerrun.aws.json" error when creating a deploy
### Byte Order Marks (BOM) in Dockerrun file
If your Dockerrun.aws.json file contains a Byte Order Marker, you may receive an "Invalid Dockerrun.aws.json" error when creating a deploy. If you create or edit the Dockerrun file using Visual Studio, and you have not modified the file encoding settings in Visual Studio, you are likely to encounter this error.

**To remove the BOM:**

* At the command line, type the following to remove the BOM:

    * (Linux/OS X) 
    
    `tail -c +4 DockerrunFile > DockerrunFileNew`
    
    Replace `DockerrunFile` with the path to your Dockerrun file, and `DockerrunFileNew` with a new name for the Dockerrun file without the BOM.

Alternatively, you can use the [dos2unix file converter](https://sourceforge.net/projects/dos2unix/) to remove the BOM from your Dockerrun files. Dos2unix is available for Windows, Linux and Mac OS.

**To remove the BOM using dos2unix:**

* At the command line, type the following:

    `dos2unix --remove-bom -n DockerrunFile DockerrunFileNew`

Replace DockerrunFile with the path to your Dockerrun file, and DockerrunFileNew with a new name for the Dockerrun file without the BOM.

---

## "AWS Error: the key pair '<keyvalue>' does not exist (code 'ValidationError')" with l0-setup

This occurs when you pass an invalid EC2 keypair to l0-setup. To fix this, follow the instructions for [creating an EC2 Key Pair](../setup/install/#part-2-create-an-access-key).

1. After you've created a new EC2 Key Pair, use [l0-setup init](reference/setup-cli/#init) to reconfigure your instance:

```
l0-setup init --aws-ssh-key-pair keypair
```

<!--
##"Back-end server is at capacity" (status code 503) error when executing Layer0 commands

The "server is at capacity" error indicates that the API server has run out of disk space. The fastest way to solve this issue is to rebuild your Layer0 API server.

**To rebuild the API server:**

1. At the command line, type the following command to force the API server to be recreated:
    * **l0-setup terraform** *Layer0Prefix* **taint aws_elastic_beanstalk.api**

2. At the command line, type the following command to re-create the API server:
    * **l0-setup apply** *Layer0Prefix*
-->
