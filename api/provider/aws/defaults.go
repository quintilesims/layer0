package aws

import "github.com/quintilesims/layer0/common/models"

const (
	DefaultInstanceSize  = "m3.medium"
	DefaultEnvironmentOS = "linux"
	DefaultServiceScale  = 1
)

const DefaultLinuxUserdataTemplate = `
#!/bin/bash
echo ECS_CLUSTER={{ .ECSEnvironmentID }} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
yum install -y aws-cli awslogs jq
aws s3 cp s3://{{ .S3Bucket }}/bootstrap/dockercfg dockercfg
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs
`

const DefaultWindowsUserdataTemplate = `
<powershell>
# Set agent env variables for the Machine context (durable)
$clusterName = "{{ .ECSEnvironmentID }}"
Write-Host Cluster name set as: $clusterName -foreground green

[Environment]::SetEnvironmentVariable("ECS_CLUSTER", $clusterName, "Machine")
[Environment]::SetEnvironmentVariable("ECS_ENABLE_TASK_IAM_ROLE", "false", "Machine")
$agentVersion = 'v1.14.0-1.windows.1'
$agentZipUri = "https://s3.amazonaws.com/amazon-ecs-agent/ecs-agent-windows-$agentVersion.zip"
$agentZipMD5Uri = "$agentZipUri.md5"

# Configure docker auth
Read-S3Object -BucketName {{ .S3Bucket }} -Key bootstrap/dockercfg -File dockercfg.json
$dockercfgContent = [IO.File]::ReadAllText("dockercfg.json")
[Environment]::SetEnvironmentVariable("ECS_ENGINE_AUTH_DATA", $dockercfgContent, "Machine")
[Environment]::SetEnvironmentVariable("ECS_ENGINE_AUTH_TYPE", "dockercfg", "Machine")

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
`

const DefaultAssumeRolePolicy = `{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ecs.amazonaws.com"
        ]
      },
      "Action": [
        "sts:AssumeRole"
      ]
    }
  ]
}`

const DefaultLBRolePolicyTemplate = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:Describe*",
                "ec2:Describe*"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
                "elasticloadbalancing:RegisterInstancesWithLoadBalancer"
            ],
            "Resource": [
                "arn:aws:elasticloadbalancing:{{ .Region }}:{{ .AccountID }}:loadbalancer/{{ .LoadBalancerID }}"
            ]
        }
    ]
}`

var DefaultLoadBalancerPort = models.Port{
	ContainerPort: 80,
	HostPort:      80,
	Protocol:      "tcp",
}

var DefaultHealthCheck = models.HealthCheck{
	Target:             "TCP:80",
	Interval:           30,
	Timeout:            5,
	HealthyThreshold:   2,
	UnhealthyThreshold: 2,
}
