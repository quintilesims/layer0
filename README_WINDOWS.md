# Layer0 Windows Support

Windows container support for Layer0 is currently alpha quality, and requires some manual steps to get working. We've documented these steps below. We expect to fully support Windows containers sometime in April, 2017.

## Instructions

- `l0 environment create windows-poc`
- Go to AWS ECS Console and record the **hashed environment name** that was created, it will be something like `l0-Layer0instancename-windows9f05c`.
- Go to the AWS EC2 Console, navigate to Auto-Scaling | Launch Configurations.
- Find the Launch Configuration that matches your hashed environment name (i.e. `l0-Layer0instancename-windows9f05c`), select it, go to Actions and select "Copy launch configuration".
- In the copied launch configuration, select a new AMI. In AMI selection, under "quick start", find the newest AMI that matches the prefix `Windows_Server-2016-English-Full-Containers`.
- Edit the Launch Configuration Details. Expand the section that reads "Advanced Details". In here, replace the "User Data" text box with the following script, replacing the `clusterName` variable with the hashed environment name from the steps above:

```
<powershell>
## The string 'windows' should be replaced with your cluster name

# Set agent env variables for the Machine context (durable)
$clusterName = "windows-cluster"
Write-Host Cluster name set as: $clusterName -foreground green

[Environment]::SetEnvironmentVariable("ECS_CLUSTER", $clusterName, "Machine")
[Environment]::SetEnvironmentVariable("ECS_ENABLE_TASK_IAM_ROLE", "false", "Machine")
$agentVersion = 'v1.14.0-1.windows.1'
$agentZipUri = "https://s3.amazonaws.com/amazon-ecs-agent/ecs-agent-windows-$agentVersion.zip"
$agentZipMD5Uri = "$agentZipUri.md5"

# Output all environment variables
Get-ChildItem -Path Env:* | Sort-Object Name

### --- Nothing user configurable after this point ---
$ecsExeDir = "$env:ProgramFiles\Amazon\ECS"
$zipFile = "$env:TEMP\ecs-agent.zip"
$md5File = "$env:TEMP\ecs-agent.zip.md5"

### Get the files from S3
Invoke-RestMethod -OutFile $zipFile -Uri $agentZipUri
Invoke-RestMethod -OutFile $md5File -Uri $agentZipMD5Uri

## MD5 Checksum
$expectedMD5 = (Get-Content $md5File)
$md5 = New-Object -TypeName System.Security.Cryptography.MD5CryptoServiceProvider
$actualMD5 = [System.BitConverter]::ToString($md5.ComputeHash([System.IO.File]::ReadAllBytes($zipFile))).replace('-', '')

if($expectedMD5 -ne $actualMD5) {
    echo "Download doesn't match hash."
    echo "Expected: $expectedMD5 - Got: $actualMD5"
    exit 1
}

## Put the executables in the executable directory.
Expand-Archive -Path $zipFile -DestinationPath $ecsExeDir -Force

## Start the agent script in the background.
$jobname = "ECS-Agent-Init"
$script =  "cd '$ecsExeDir'; .\amazon-ecs-agent.ps1"
$repeat = (New-TimeSpan -Minutes 1)

$jobpath = $env:LOCALAPPDATA + "\Microsoft\Windows\PowerShell\ScheduledJobs\$jobname\ScheduledJobDefinition.xml"
if($(Test-Path -Path $jobpath)) {
  echo "Job definition already present"
  exit 0

}

$scriptblock = [scriptblock]::Create("$script")
$trigger = New-JobTrigger -At (Get-Date).Date -RepeatIndefinitely -RepetitionInterval $repeat -Once
$options = New-ScheduledJobOption -RunElevated -ContinueIfGoingOnBattery -StartIfOnBattery
Register-ScheduledJob -Name $jobname -ScriptBlock $scriptblock -Trigger $trigger -ScheduledJobOption $options -RunNow
Add-JobTrigger -Name $jobname -Trigger (New-JobTrigger -AtStartup -RandomDelay 00:1:00)
</powershell>
<persist>true</persist>
```

- Under IP Address Type, select `Assign a public IP address to every instance.`. We will use this public IP (which is obtained when the instance is created) for communicating with an existing Layer0 environment and/or RDP.
- Go to the 'Add Storage' tab; select an appropriate disk size. Since we aren't sure how big to make these instances yet, let's default to 128gb for now.
- Go the 'Security Group' tab; select a single default security group. This group should already exist and would be the same name as your hashed environment name from above steps with the added suffix `-env`.
- Proceed to Review and complete the creation of the Launch Configuration.
- In the EC2 Console, Navigate to Auto-Scaling | Auto Scaling Groups.
- Select the Auto Scaling Group that matches the hashed environment name from above, go to Actions and Select "Edit".
- In the edit panel, select the Launch Configuration that you just created in the `Launch Configuration` dropdown and click on "Save".

At the point, new instances created for the environment specified in step 1 will be created using our new Launch Configuration settings. During our initial testing, a new Windows instance takes a bit longer than a Linux instance to come up (~5 mins), and registering with the ECS cluster takes about 3 minutes.

You should be able to create new Windows [deploys](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/windows_task_definitions.html) and services using the `l0` binary from this point forth for the environment you specified. Let us know if you encounter issues!

## Additional Caveats

- A Layer0 environment can _only_ be a Windows or Linux environment --- co-mingling container types in a single environment is not allowed.
- Windows networking is different than Linux so do not rely on things like `localhost` for local resolution.
- [AWS Windows / ECS Caveats ](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_Windows.html)

## Windows Deploy (Task Definitions) Reference
- http://docs.aws.amazon.com/AmazonECS/latest/developerguide/windows_task_definitions.html
