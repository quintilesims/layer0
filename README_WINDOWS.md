# Layer0 Windows Support

Windows container support for Layer0 is currently experimental, and requires some manual steps to get working. We've documented these steps below. We expect to fully support Windows containers sometime in April, 2017.

## Requirements

- A Layer0 Instance with the `l0` binary [properly configured](https://quintilesims.github.io/layer0/setup/install/).
- The suplemental file `win-user-data.ps1`, which is located in the `/scripts` directory of the [Layer0 source repository](https://github.com/quintilesims/layer0).

## Instructions

- `l0 environment create --user-data win-user-data.ps1 windows-poc`

- Go to AWS ECS Console and record the **hashed environment name** that was created, it will be something like `l0-Layer0instancename-windows9f05c`.

- Go to the AWS EC2 Console, navigate to Auto-Scaling | Launch Configurations.

- Find the Launch Configuration that matches your hashed environment name (i.e. `l0-Layer0instancename-windows9f05c`), select it, go to Actions and select "Copy launch configuration".

- A newly copied launch configuration will be made that matches the name of your hashed environment name with a suffix "Copy" (you can rename as you see fit, but this is the default). Select the copied launch configuration, and in the edit panel, select a new AMI. Under AMI selection, select the "Quick Start" column, and find the newest AMI that matches the prefix `Windows_Server-2016-English-Full-Containers`.

- Under IP Address Type, select `Assign a public IP address to every instance.`. We will use this public IP (which is obtained when the instance is created) for communicating with an existing Layer0 environment and/or RDP.

- Go to the 'Add Storage' tab; select an appropriate disk size. Since we aren't sure how big to make these instances yet, let's use AWS' walkthrough default of 200gb.

- Go to the 'Security Group' tab; select a single default security group. This group should already exist and has same name as your hashed environment name from above steps with the added suffix `-env`.

- Proceed to Review and complete the creation of the Launch Configuration.

- In the EC2 Console, Navigate to Auto-Scaling | Auto Scaling Groups.

- Select the Auto Scaling Group that matches the hashed environment name from above, go to Actions and Select "Edit".

- In the edit panel, select the Launch Configuration that you just created in the `Launch Configuration` dropdown and click on "Save".

At the point, new instances created for the environment specified in step 1 will be created using our new Launch Configuration settings. During our initial testing, a new Windows instance takes a bit longer than a Linux instance to come up (~5 mins), and registering with the ECS cluster takes about 3 minutes.

You should be able to create new Windows [deploys](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/windows_task_definitions.html) and services using the `l0` binary from this point forth for the environment you specified. Let us know if you encounter issues!

## Additional Caveats

- A Layer0 environment can _only_ be a Windows or Linux environment --- co-mingling container types in a single environment is not allowed.
- Windows networking is different than Linux so don't rely on `localhost` for local name resolution.
- [AWS Windows / ECS Caveats ](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_Windows.html)

## Windows Deploy (Task Definitions) Reference
- http://docs.aws.amazon.com/AmazonECS/latest/developerguide/windows_task_definitions.html
