# Install and Configure Layer0

## Prerequisites

Before you can install and configure Layer0, you must obtain the following:

* **An AWS account.**
* **An EC2 Key Pair.** This key pair allows you to access the EC2 instances running your Services using SSH. If you have already created a key pair, you can use it for this process. Otherwise, follow the [instructions at aws.amazon.com](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html#having-ec2-create-your-key-pair) to create a new key pair. Make a note of the name that you selected when creating the key pair.

## Part 1: Download and extract Layer0

1. In the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the zip file to a directory on your computer.
2. Add both the directory that contains the **l0** application, as well as the entire l0-setup directory, to your system path. The l0-setup directory contains files that are necessary for the **l0-setup** application to work properly, so you must add the entire directory to your path.<br />For more information about adding directories to your system path, see the following resources:
	* (Windows): [How to Edit Your System PATH for Easy Command Line Access in Windows](http://www.howtogeek.com/118594/how-to-edit-your-system-path-for-easy-command-line-access/)
	* (Linux/macOS): [Adding a Directory to the Path](http://www.troubleshooters.com/linux/prepostpath.htm)

## Part 2: Create an Access Key
This step will create an Identity & Access Management (IAM) access key from your AWS account. You will use the credentials created in this section when installing, updating, or removing Layer0 resources.

**To create an Access Key:**

1. In a web browser, login to the [AWS Console](http://console.aws.amazon.com/).

2. Under **Security and Identity**, click **Identity and Access Management**.

3. Click **Groups**, and then click **Administrators**. <div class="admonition note"><p class="admonition-title">Note</p><br /><p>If the **Administrators** group does not already exist, complete the following steps: <ol><li>Click **Create New Group**. Name the new group "Administrators", and then click **Next Step**.</li><li>Click **AdministratorAccess** to attach the Administrator policy to your new group.</li><li>Click **Next Step**, and then click **Create Group**.</li></ul></p></div>

4. Click **Users**.

5. Click **Create New Users** and enter a unique user name you will use for Layer0. This user name can be used for multiple Layer0 installations. Check the box next to **Generate an Access Key for each user**, and then click **Create**.

6. Once your user account has been created, click **Download Credentials** to save your access key to a CSV file.

7. In the Users list, click the user account you just created. Under **User Actions**, click **Add User to Groups**.

8. Select the group **Administrators** and click **Add to Groups**. This will make your newly created user an administrator for your AWS account, so be sure to keep your security credentials safe!

## Part 3: Configure your Layer0
Now that you have downloaded Layer0 and configured your AWS instance, you can create your Layer0.

**To configure Layer0:**

1. At the command prompt, navigate to the **l0-setup** subdirectory in the folder in which you extracted the Layer0 files.
2. Type the following command, replacing ``[prefix]`` with a unique name for your Layer0: ```l0-setup apply [prefix]```
3. When prompted, enter the following information:
	* **AWS Access Key ID**: The access key ID contained in the credential file that you downloaded in step 6 of the previous section.
	* **AWS Secret Access Key**: The secret access key contained in the credential file that you downloaded in step 6 of the previous section.
	* **Key Pair**: The name of the key pair that you created in the Prerequisites section.
The first time you run the ```apply``` command, it may take around 15 minutes to complete. If the ```apply``` command fails to complete successfully, it is safe to run it again until it succeeds.

## (Optional) Part 4: Configure the dockercfg file to use a private Docker registry

!!! note
	The procedures in this section are optional, but are highly recommended for production use.

When you run the ```l0-setup apply``` command for the first time, a blank dockercfg file will be created in a folder in your layer0 directory corresponding to the name of your Layer0 prefix. You can modify this file to use a private Docker registry by completing the following steps:

<ol>
	<li>To add private registry authentication, modify this file to include the authentication information in the following format:
<pre class="code"><code>{
  "https://index.docker.io/v1/": {
    "username": "my_name",
    "password": "my_password",
    "email": "email@example.com"
  }
}</code></pre>
	</li>
	<li>Save the modified dockercfg file, and then run the following command: <code>l0-setup apply [prefix]</code></li>
</ol>

## Part 5: Configure the environment variables
Once the ```apply``` command has run successfully, you can configure the Layer0 environment variables using the ```endpoint``` command.

To view the environment variables for your Layer0 and apply them to your shell, type the following command, replacing ```[prefix]``` with the name of the Layer0 prefix you created in Part 3:

* (Windows PowerShell): ```l0-setup endpoint --insecure --powershell [prefix] | Out-String | Invoke-Expression```
* (Linux/macOS): ```eval "$(l0-setup endpoint --insecure [prefix])"```

## (Optional) Part 6: Using a custom certificate

!!! note
	The procedures in this section are optional, but are highly recommended for production use.

Layer0 uses a self-signed certificate to run the API. In some cases, you may want to use a custom certificate instead.

**To use a custom certificate:**

<ol>
	<li>In a text editor, open the file in your layer0 directory named elb.tf.template. Insert the following line at the end of the file:
<pre class="code"><code>
...
	ssl_certificate_id = "[ARN of your SSL cert]"
...
</pre></code></li>
	<li>Save the modified elb.tf.template file, and then run the following command: <code>l0-setup apply [prefix]</code></li>
</ol>
