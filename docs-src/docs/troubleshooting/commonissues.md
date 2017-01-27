#Common issues and their solutions

##"Connection refused" error when executing Layer0 commands

When executing commands using the Layer0 CLI, you may see the following error message: "Get http://localhost:9090/_command_/: dial tcp 127.0.0.1:9090: connection refused", where _command_ is the Layer0 command you are trying to execute.

This error indicates that your Layer0 environment variables have not been set for the current session. See the ["Configure environment variables" section](http://localhost:8000/setup/install/#part-4-configure-environment-variables) of the Layer0 installation guide for instructions for setting up your environment variables.

---

##"Invalid Dockerrun.aws.json" error when creating a deploy
###Byte Order Marks (BOM) in Dockerrun file
If your Dockerrun.aws.json file contains a Byte Order Marker, you may receive an "Invalid Dockerrun.aws.json" error when creating a deploy. If you create or edit the Dockerrun file using Visual Studio, and you have not modified the file encoding settings in Visual Studio, you are likely to encounter this error.

**To remove the BOM:**

* At the command line, type the following to remove the BOM:

    * (Linux/OS X) **tail -c +4** _DockerrunFile_ **>** _DockerrunFileNew_
    <br /><br />Replace _DockerrunFile_ with the path to your Dockerrun file, and _DockerrunFileNew_ with a new name for the Dockerrun file without the BOM.

Alternatively, you can use the [dos2unix file converter](https://sourceforge.net/projects/dos2unix/) to remove the BOM from your Dockerrun files. Dos2unix is available for Windows, Linux and Mac OS.

**To remove the BOM using dos2unix:**

* At the command line, type the following:

    * **dos2unix --remove-bom -n** _DockerrunFile_ _DockerrunFileNew_
    <br /><br />Replace _DockerrunFile_ with the path to your Dockerrun file, and _DockerrunFileNew_ with a new name for the Dockerrun file without the BOM.

---

##"AWS Error: the key pair '<keyvalue>' does not exist (code 'ValidationError')" with l0-setup

This occurs when you pass a non-existent EC2 keypair to l0-setup. To fix this, follow the instructions for [creating an EC2 Key Pair](/install/#part-2-create-an-access-key).

1. After you've created a new EC2 Key Pair, run the following command:
<ul>
  <li class="command">**l0-setup plan** *prefix* **-var key_pair**=*keypair*</li>
</ul>

<!--
##"Back-end server is at capacity" (status code 503) error when executing Layer0 commands

The "server is at capacity" error indicates that the API server has run out of disk space. The fastest way to solve this issue is to rebuild your Layer0 API server.

**To rebuild the API server:**

1. At the command line, type the following command to force the API server to be recreated:
    * **l0-setup terraform** *Layer0Prefix* **taint aws_elastic_beanstalk.api**

2. At the command line, type the following command to re-create the API server:
    * **l0-setup apply** *Layer0Prefix*
-->
