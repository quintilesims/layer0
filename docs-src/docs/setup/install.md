# Create a new Layer0 Instance

## Prerequisites

Before you can install and configure Layer0, you must obtain the following:

* **An AWS account.**

* **An EC2 Key Pair.**
This key pair allows you to access the EC2 instances running your Services using SSH.
If you have already created a key pair, you can use it for this process.
Otherwise, follow the [instructions at aws.amazon.com](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#having-ec2-create-your-key-pair) to create a new key pair.
Make a note of the name that you selected when creating the key pair.

* **Terraform v0.9.4+**
We use Terraform to create the resources that Layer0 needs.
If you're unfamiliar with Terraform, you may want to check out our [introduction](/reference/terraform_introduction).
If you're ready to install Terraform, there are instructions in the [Terraform documentation](https://www.terraform.io/intro/getting-started/install.html).

## Part 1: Download and extract Layer0

1. In the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the zip file to a directory on your computer.
2. (Optional) Place the **l0** and **l0-setup** binaries into your system path. 
For more information about adding directories to your system path, see the following resources:
	* (Windows): [How to Edit Your System PATH for Easy Command Line Access in Windows](http://www.howtogeek.com/118594/how-to-edit-your-system-path-for-easy-command-line-access/)
	* (Linux/macOS): [Adding a Directory to the Path](http://www.troubleshooters.com/linux/prepostpath.htm)

## Part 2: Create an Access Key
This step will create an Identity & Access Management (IAM) access key for your AWS account. 
You will use the credentials created in this section when creating, updating, or removing Layer0 instances.

**To create an Access Key:**

1. In a web browser, login to the [AWS Console](http://console.aws.amazon.com/).

2. Under **Security and Identity**, click **Identity and Access Management**.

3. Click **Groups**, and then click **Administrators**. <div class="admonition note"><p class="admonition-title">Note</p><br /><p>If the **Administrators** group does not already exist, complete the following steps: <ol><li>Click **Create New Group**. Name the new group "Administrators", and then click **Next Step**.</li><li>Click **AdministratorAccess** to attach the Administrator policy to your new group.</li><li>Click **Next Step**, and then click **Create Group**.</li></ul></p></div>

4. Click **Users**.

5. Click **Create New Users** and enter a unique user name you will use for Layer0. This user name can be used for multiple Layer0 installations. Check the box next to **Generate an Access Key for each user**, and then click **Create**.

6. Once your user account has been created, click **Download Credentials** to save your access key to a CSV file.

7. In the Users list, click the user account you just created. Under **User Actions**, click **Add User to Groups**.

8. Select the group **Administrators** and click **Add to Groups**. This will make your newly created user an administrator for your AWS account, so be sure to keep your security credentials safe!

## Part 3: Create a new Layer0 Instance
Now that you have downloaded Layer0 and configured your AWS account, you can create your Layer0 instance.
From a command prompt, run the following (replacing `<instance_name>` with a name for your Layer0 instance):
```
$ l0-setup init <instance_name>
```

This command will prompt you for many different inputs. 
Enter the required values for **AWS Access Key**, **AWS Secret Key**, and **AWS SSH Key** as they come up.
All remaining inputs are optional and can be set to their default by pressing enter.

```
...
AWS Access Key: The access_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Access Key ID portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources.

[current: <none>]
Please enter a value and press 'enter'.
        Input: ABC123xzy

AWS Secret Key: The secret_key input variable is used to provision the AWS resources
required for Layer0. This corresponds to the Secret Access Key portion of an AWS Access Key.
It is recommended this key has the 'AdministratorAccess' policy. Note that Layer0 will
only use this key for 'l0-setup' commands associated with this Layer0 instance; the
Layer0 API will use its own key with limited permissions to provision AWS resources.

[current: <none>]
Please enter a value and press 'enter'.
        Input: ZXY987cba

AWS SSH Key Pair: The ssh_key_pair input variable specifies the name of the
ssh key pair to include in EC2 instances provisioned by Layer0. This key pair must
already exist in the AWS account. The names of existing key pairs can be found
in the EC2 dashboard. Note that changing this value will not effect instances
that have already been provisioned.

[current: <none>]
Please enter a value and press 'enter'.
        Input: mySSHKey
...
```

Once the **init** command has successfully completed, you're ready to actually create the resources needed to use Layer0.
Run the following command (again, replace `<instance_name>` with the name you've chosen for your Layer0 instance):

```
l0-setup apply <instance_name>
```

The first time you run the **apply** command, it may take around 5 minutes to complete. 
This command is idempotent; it is safe to run multiple times if it fails the first.

It's a good idea to run the **push** command (`l0-setup push <instance_name>`) after **apply** commands complete. 
This will send a backup of your Layer0 instance's configuration and state to S3. 
These files can be grabbed later using the **pull** command (`l0-setup pull <instance_name>`). 

!!! note "Using a Private Docker Registry"
    **The procedures in this section are optional, but are highly recommended for production use.**

If you require authentication to a private Docker registry, you will need a Docker configuration file present on your machine with access to private repositories (typically located at `~/.docker/config.json`). 

If you don't have a config file yet, you can generate one by running `docker login [registry-address]`. 
A configuration file will be generated at `~/.docker/config.json`.

To add this authentication to your Layer0 instance, run:
```
$ l0-setup init --docker-path=<path/to/config.json> <instance_name>
```

This will add a rendered file into your Layer0 instance's directory at `~/.layer0/<instance_name>/dockercfg.json`.

You can modify a Layer0 instance's `dockercfg.json` file and re-run the **apply** command (`l0-setup apply <instance_name>`) to make changes to your authentication. 
Note that any EC2 instances created prior to changing your `dockercfg.json` file will need to be manually terminated since they only grab the authentication file during instance creation. 
Terminated EC2 instances will be automatically re-created by autoscaling.


!!! note "Using an Existing VPC"
    **The procedures in this section must be followed to properly install Layer0 into an existing VPC**

By default, l0-setup creates a new VPC to place resources. 
However, l0-setup can place resources in an existing VPC if it meets the following conditions:

* Has access to the public internet (through a NAT instance or gateway)
* Has at least 1 public and 1 private subnet
* The public and private subnets have the tag `Tier: Public` or `Tier: Private`, respectively.
For information on how to tag AWS resources, please visit the [AWS documentation](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html). 

Once you are sure the existing VPC satisfies these requirements, run the **init** command, 
placing the VPC ID when prompted:
```
$ l0-setup init <instance_name>
...
VPC ID (optional): The vpc_id input variable specifies an existing AWS VPC to provision
the AWS resources required for Layer0. If no input is specified, a new VPC will be
created for you. Existing VPCs must satisfy the following constraints:

    - Have access to the public internet (through a NAT instance or gateway)
    - Have at least 1 public and 1 private subnet
    - Each subnet must be tagged with ["Tier": "Private"] or ["Tier": "Public"]

Note that changing this value will destroy and recreate any existing resources.

[current: ]
Please enter a new value, or press 'enter' to keep the current value.
        Input: vpc123
```

Once the command has completed, it is safe to run **apply** to provision the resources. 


## Part 4: Connect to a Layer0 Instance
Once the **apply** command has run successfully, you can configure the environment variables needed to connect to the Layer0 API using the **endpoint** command.

```
$ l0-setup endpoint --insecure <instance_name>
export LAYER0_API_ENDPOINT="https://l0-instance_name-api-123456.us-west-2.elb.amazonaws.com"
export LAYER0_AUTH_TOKEN="abcDEFG123"
export LAYER0_SKIP_SSL_VERIFY="1"
export LAYER0_SKIP_VERSION_VERIFY="1"
```

The **--insecure** flag shows configurations that bypass SSL and version verifications. 
This is required as the Layer0 API created uses a self-signed certificate by default.
These settings are **not** recommended for production use!

The **endpoint** command supports a `--syntax` option, which can be used to turn configuration into a single line:

* Bash (default) - `$ eval "$(l0-setup endpoint --insecure <instance_name>)"`
* Powershell - `$ l0-setup endpoint --insecure --syntax=powershell <instance_name> | Out-String | Invoke-Expression`
