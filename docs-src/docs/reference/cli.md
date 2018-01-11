# Layer0 CLI Reference

##Global options

The `l0` application is designed to be used with one of several commands: [admin](#admin), [deploy](#deploy), [environment](#environment), [job](#job), [loadbalancer](#loadbalancer), [service](#service), and [task](#task). These commands are detailed in the sections below. There are, however, some global parameters that you may specify whenever using `l0`.

####Usage
```
l0 [global options] command subcommand [subcommand options] params
```

####Global options
* `-o [text|json], --output [text|json]` - Specify the format of Layer0 outputs. By default, Layer0 outputs unformatted text; by issuing the `--output json` option, you can force `l0` to output JSON-formatted text.
* `-t value, --timeout value` - Specify the timeout for running `l0` commands. Values can be in h, m, s, or ms.
* `-d, --debug` - Print debug statements
* `--version` - Display the version number of the `l0` application.

---

##Admin
The `admin` command is used to manage the Layer0 API server. This command is used with the following subcommands: [debug](#admin-debug), [sql](#admin-sql), and [version](#admin-version).

### admin debug
Use the `debug` subcommand to view the running version of your Layer0 API server and CLI.

#### Usage
```
l0 admin debug
```

### admin sql
Use the `sql` subcommand to initialize the Layer0 API database.

#### Usage
```
l0 admin sql
```

####Additional information
The `sql` subcommand is automatically executed during the Layer0 installation process; we recommend that you do not use this subcommand unless specifically directed to do so.

### admin version
Use the `version` subcommand to display the current version of the Layer0 API.

#### Usage
```
l0 admin version 
```

---

##Deploy
Deploys are ECS Task Definitions. They are configuration files that detail how to deploy your application.
The `deploy` command is used to manage Layer0 environments. This command is used with the following subcommands: [create](#deploy-create), [delete](#deploy-delete), [get](#deploy-get), and [list](#deploy-list).

### deploy create
Use the `create` subcommand to upload a Docker task definition into Layer0. 

#### Usage
```
l0 deploy create dockerPath deployName
```

####Required parameters
* `dockerPath` - The path to the Docker task definition that you want to upload.
* `deployName` - A name for the deploy.

####Additional information
If `deployName` exactly matches the name of an existing Layer0 deploy, then the version number of that deploy will increase by 1, and the task definition you specified will replace the task definition specified in the previous version.

If you use Visual Studio to modify or create your Dockerrun file, you may see an "Invalid Dockerrun.aws.json" error. This error is caused by the default encoding used by Visual Studio. See the ["Common issues" page](http://localhost:8000/troubleshooting/commonissues/#invalid-dockerrunawsjson-error-when-creating-a-deploy) for steps to resolve this issue.

Deploys created through Layer0 are rendered with a `logConfiguration` section for each container.
If a `logConfiguration` section already exists, no changes are made to the section.
The additional section enables logs from each container to be sent to the the Layer0 log group.
This is where logs are looked up during `l0 <entity> logs` commands.
The added `logConfiguration` section uses the following template:

```
"logConfiguration": {
	"logDriver": "awslogs",
		"options": {
			"awslogs-group": "l0-<prefix>",
			"awslogs-region": "<region>",
			"awslogs-stream-prefix": "l0"
		}
	}
}
```

### deploy delete
Use the `delete` subcommand to delete a version of a Layer0 deploy.

#### Usage
```
l0 deploy delete deployName
```

####Required parameters
* `deployName` - The name of the Layer0 deploy you want to delete.

### deploy get
Use the `get` subcommand to view information about an existing Layer0 deploy.

#### Usage
```
l0 deploy get deployName
```

####Required parameters
* `deployName` - The name of the Layer0 deploy for which you want to view additional information.

####Additional information
The `get` subcommand supports wildcard matching: `l0 deploy get dep*` would return all deploys beginning with `dep`.

### deploy list
Use the `list` subcommand to view a list of deploys in your instance of Layer0.

#### Usage
```
l0 deploy list
```

---

## Environment
Layer0 environments allow you to isolate services and load balancers for specific applications.
The `environment` command is used to manage Layer0 environments. This command is used with the following subcommands: [create](#environment-create), [delete](#environment-delete), [get](#environment-get), [list](#environment-list), and [setmincount](#environment-setmincount).

### environment create
Use the `create` subcommand to create a new Layer0 environment.

#### Usage
```
l0 environment create [--size size | --min-count mincount | 
    --user-data path | --os os | --ami amiID] environmentName
```

####Required parameters
* `environmentName` - A name for the environment.

####Optional arguments
* `--size size` - The instance size of the EC2 instances to create in your environment (default: m3.medium).
* `--min-count mincount` - The minimum number of EC2 instances allowed in the environment's autoscaling group (default: 0).
* `--user-data path` - The user data template file to use for the environment's autoscaling group.
* `--os os` - The operating system used in the environment. Options are "linux" or "windows" (default: linux). More information on windows environments is documented below.
* `ami amiID` - A custom EC2 AMI ID to use in the environment. If not specified, Layer0 will use its default AMI ID for the specified operating system.

The user data template can be used to add custom configuration to your Layer0 environment. They are usually scripts that are executed at instance launch time to ensure an EC2 instance is in the correct state after the provisioning process finishes.
Layer0 uses [Go Templates](https://golang.org/pkg/text/template) to render user data.
Currently, two variables are passed into the template: **ECSEnvironmentID** and **S3Bucket**.

!!! danger
    Please review the [ECS Tutorial](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_container_instance.html)
    to better understand how to write a user data template, and use at your own risk!


**Linux Environments**: The default Layer0 user data template is:

``` bash 
#!/bin/bash
echo ECS_CLUSTER={{ .ECSEnvironmentID }} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
yum install -y aws-cli awslogs jq
aws s3 cp s3://{{ .S3Bucket }}/bootstrap/dockercfg dockercfg
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs
```

**Windows Environments**: The default Layer0 user data template is:
``` powershell
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
```


!!! note "Windows Environments"
        Windows containers are still in beta. 
You can view the documented caveats with ECS [here](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_Windows.html#windows_caveats).
When creating Windows environments in Layer0, the root volume sizes for instances are 200GiB to accommodate the large size of the containers.  
It can take as long as 45 minutes for a new windows container to come online. 

### environment delete
Use the `delete` subcommand to delete an existing Layer0 environment.

#### Usage
```
l0 environment delete [--wait] environmentName
```

####Required parameters
* `environmentName` - The name of the Layer0 environment that you want to delete.

####Optional arguments
* `--wait` - Wait until the deletion is complete before exiting.

####Additional information
This operation performs several tasks asynchronously. When run without the `--wait` option, this operation will most likely exit before all of these tasks are complete; when run with the `--wait` option, this operation will only exit once these tasks have completed.

### environment get
Use the `get` subcommand to display information about an existing Layer0 environment.

#### Usage
```
l0 environment get environmentName
```

####Required parameters
* `environmentName` - The name of the Layer0 environment for which you want to view additional information.

####Additional information
The `get` subcommand supports wildcard matching: `l0 environment get test*` would return all environments beginning with `test`.

### environment list
Use the `list` subcommand to display a list of environments in your instance of Layer0.

#### Usage
```
l0 environment list
```

### environment setmincount
Use the `setmincount` subcommand to set the minimum number of EC2 instances allowed the environment's autoscaling group.

####Usage
```
l0 environment setmincount environmentName count
```

####Required parameters
* `environmentName` - The name of the Layer0 environment that you want to adjust.
* `count` - The minimum number of instances allowed in the environment's autoscaling group.

### environment link
Use the `link` subcommand to link two environments together. 
When environments are linked, services inside the environments are allowed to communicate with each other as if they were in the same environment. 
This link is bidirectional. 
This command is idempotent; it will succeed even if the two specified environments are already linked.

####Usage
```
l0 environment link sourceEnvironmentName destEnvironmentName
```

####Required parameters
* `sourceEnvironmentName` - The name of the source environment to link.
* `destEnvironmentName` - The name of the destination environment to link.

### environment unlink
Use the `unlink` subcommand to remove the link between two environments.
This command is idempotent; it will succeed even if the link does not exist.

####Usage
```
l0 environment unlink sourceEnvironmentName destEnvironmentName
```

####Required parameters
* `sourceEnvironmentName` - The name of the source environment to unlink.
* `destEnvironmentName` - The name of the destination environment to unlink.

---

##Job
A Job is a long-running unit of work performed on behalf of the Layer0 API.
Jobs are executed as Layer0 tasks that run in the **api** environment.
The `job` command is used with the following subcommands: [logs](#job-logs), [delete](#job-delete), [get](#job-get), and [list](#job-list).

### job logs
Use the `logs` subcommand to display the logs from a Layer0 job that is currently running.

####Usage
```
l0 job logs [--start MM/DD HH:MM | --end MM/DD HH:MM | --tail=N] jobName
```

####Required parameters
* `jobName` - The name of the Layer0 job for which you want to view logs.

####Optional arguments
* `--start MM/DD HH:MM` - The start of the time range to fetch logs.
* `--end MM/DD HH:MM` - The end of the time range to fetch logs.
* `--tail=N` - Display only the last `N` lines of the log.

###job delete
Use the `delete` subcommand to delete an existing job.

####Usage
```
l0 job delete jobName
```

####Required parameters
* `jobName` - The name of the job that you want to delete.

###job get
Use the `get` subcommand to display information about an existing Layer0 job.

####Usage
```
l0 job get jobName
```

####Required parameters
* `jobName` - The name of an existing Layer0 job to display.

####Additional information
The `get` subcommand supports wildcard matching: `l0 job get 2a55*` would return all jobs beginning with `2a55`.

###job list
Use the `list` subcommand to display information about all of the existing jobs in an instance of Layer0.

####Usage
```
l0 job list
```

---

##Loadbalancer
A load balancer is a component of a Layer0 environment. Load balancers listen for traffic on certain ports, and then forward that traffic to Layer0 [services](#service). The `loadbalancer` command is used with the following subcommands: [create](#loadbalancer-create), [delete](#loadbalancer-delete), [addport](#loadbalancer-addport), [dropport](#loadbalancer-dropport), [get](#loadbalancer-get), [list](#loadbalancer-list), and [healthcheck](#loadbalancer-healthcheck).

###loadbalancer create
Use the `create` subcommand to create a new load balancer.

####Usage
```
l0 loadbalancer create [--port port ... | --certificate certifiateName | 
    --private | --healthcheck-target target | --healthcheck-interval interval | 
    --healthcheck-timeout timeout | --healthcheck-healthy-threshold healthyThreshold | 
    --healthcheck-unhealthy-threshold unhealthyThreshold] environmentName loadBalancerName
```

####Required parameters
* `environmentName` - The name of the existing Layer0 environment in which you want to create the load balancer.
* `loadBalancerName` - A name for the load balancer you are creating.

####Optional arguments
* `--port port ...` - The port configuration for the listener of the load balancer. Valid pattern is `hostPort:containerPort/protocol`. Multiple ports can be specified using `--port port1 --port port2 ...` (default: `80/80:TCP`).
    * `hostPort` - The port that the load balancer will listen for traffic on.
    * `containerPort` - The port that the load balancer will forward traffic to.
    * `protocol` - The protocol to use when forwarding traffic (acceptable values: TCP, SSL, HTTP, and HTTPS).
* `--certificate certificateName` - The name of an existing Layer0 certificate. You must include this option if you are using an HTTPS port configuration.
* `--private` - When you use this option, the load balancer will only be accessible from within the Layer0 environment.
* `--healthcheck-target target` - The target of the check. Valid pattern is `PROTOCOL:PORT/PATH` (default: `"TCP:80"`). 
    * If `PROTOCOL` is `HTTP` or `HTTPS`, both `PORT` and `PATH` are required. Example: `HTTP:80/admin/healthcheck`. 
    * If `PROTOCOL` is `TCP` or `SSL`, `PORT` is required and `PATH` is not used. Example: `TCP:80`
* `--healthcheck-interval interval` - The interval between checks (default: `30`).
* `--healthcheck-timeout timeout` - The length of time before the check times out (default: `5`).
* `--healthcheck-healthy-threshold healthyThreshold` - The number of checks before the instance is declared healthy (default: `2`).
* `--healthcheck-unhealthy-threshold unhealthyThreshold` - The number of checks before the instance is declared unhealthy (default: `2`).

!!! info "Ports and Health Checks"
    When both the `--port` and the `--healthcheck-target` options are omitted, Layer0 configures the load balancer with some default values: `80:80/TCP` for ports and `TCP:80` for healthcheck target.
    These default values together create a load balancer configured with a simple but functioning health check, opening up a set of ports that allows traffic to the target of the healthcheck.
    (`--healthcheck-target TCP:80` tells the load balancer to ping its services at port 80 to determine their status, and `--port 80:80/TCP` configures a security group to allow traffic to pass between port 80 of the load balancer and port 80 of its services)

    When creating a load balancer with non-default configurations for either `--port` or `--healthcheck-target`, make sure that a valid `--port` and `--healthcheck-target` pairing is also created.

###loadbalancer delete
Use the `delete` subcommand to delete an existing load balancer.

####Usage
```
l0 loadbalancer delete [--wait] loadBalancerName
```

####Required parameters
* `loadBalancerName` - The name of the load balancer that you want to delete.

####Optional arguments
* `--wait` - Wait until the deletion is complete before exiting.

####Additional information
In order to delete a load balancer that is already attached to a service, you must first delete the service that uses the load balancer.

This operation performs several tasks asynchronously. When run without the `--wait` option, this operation will most likely exit before all of these tasks are complete; when run with the `--wait` option, this operation will only exit once these tasks have completed
.
###loadbalancer addport
Use the `addport` subcommand to add a new port configuration to an existing Layer0 load balancer.

####Usage
```
l0 loadbalancer addport [--certificate certificateName] loadBalancerName port
```

####Required parameters
* `loadBalancerName` - The name of an existing Layer0 load balancer in which you want to add the port configuration.
* `port` - The port configuration for the listener of the load balancer. Valid pattern is `hostPort:containerPort/protocol`.
    * `hostPort` - The port that the load balancer will listen for traffic on.
    * `containerPort` - The port that the load balancer will forward traffic to.
    * `protocol` - The protocol to use when forwarding traffic (acceptable values: TCP, SSL, HTTP, and HTTPS).

####Optional arguments
* `--certificate certificateName` - The name of an existing Layer0 certificate. You must include this option if you are using an HTTPS port configuration.

####Additional information
The port configuration you specify must not already be in use by the load balancer you specify.

###loadbalancer dropport
Use the `dropport` subcommand to remove a port configuration from an existing Layer0 load balancer.

####Usage
```
l0 loadbalancer dropport loadBalancerName hostPort
```

####Required parameters
* `loadBalancerName`- The name of an existing Layer0 load balancer from which you want to remove the port configuration.
* `hostPort`- The host port to remove from the load balancer.

###loadbalancer get
Use the `get` subcommand to display information about an existing Layer0 load balancer.

####Usage
```
l0 loadbalancer get [environmentName:]loadBalancerName
```

####Required parameters
* `[environmentName:]loadBalancerName` - The name of an existing Layer0 load balancer. You can optionally provide the Layer0 environment (`environmentName`) associated with the Load Balancer

####Additional information
The `get` subcommand supports wildcard matching: `l0 loadbalancer get entrypoint*` would return all jobs beginning with `entrypoint`.

###loadbalancer list
Use the `list` subcommand to display information about all of the existing load balancers in an instance of Layer0.

####Usage
```
l0 loadbalancer list
```

###loadbalancer healthcheck
Use the `healthcheck` subcommand to display information about or update the configuration of a load balancer's health check.

####Usage
```
l0 loadbalancer healthcheck [--set-target target | --set-interval interval | 
    --set-timeout timeout | --set-healthy-threshold healthyThreshold | 
    --set-unhealthy-threshold unhealthyThreshold] loadbalancerName
```

####Required parameters
* `loadBalancerName` - The name of the existing Layer0 load balancer you are modifying.

####Optional arguments
* `--set-target target` - The target of the check. Valid pattern is `PROTOCOL:PORT/PATH`.
    * If `PROTOCOL` is `HTTP` or `HTTPS`, both `PORT` and `PATH` are required. Example: `HTTP:80/admin/healthcheck`.
    * If `PROTOCOL` is `TCP` or `SSL`, `PORT` is required and `PATH` is not used. Example: `TCP:80`
* `--set-interval interval` - The interval between health checks.
* `--set-timeout timeout` - The length of time in seconds before the health check times out.
* `--set-healthy-threshold healthyThreshold` - The number of checks before the instance is declared healthy.
* `--set-unhealthy-threshold unhealthyThreshold` - The number of checks before the instance is declared unhealthy.

####Additional information
Calling the subcommand without flags will display the current configuration of the load balancer's health check. Setting any of the flags will update the corresponding field in the health check, and all omitted flags will leave the corresponding fields unchanged.

---

## Service
A service is a component of a Layer0 environment. The purpose of a service is to execute a Docker image specified in a [deploy](#deploy). In order to create a service, you must first create an [environment](#environment) and a [deploy](#deploy); in most cases, you should also create a [load balancer](#loadbalancer) before creating the service.

The `service` command is used with the following subcommands: [create](#service-create), [delete](#service-delete), [get](#service-get), [update](#service-update), [list](#service-list), [logs](#service-logs), and [scale](#service-scale).

###service create
Use the `create` subcommand to create a Layer0 service.

####Usage
```
l0 service create [--loadbalancer [environmentName:]loadBalancerName | 
    --no-logs] environmentName serviceName deployName[:deployVersion]
```

####Required parameters
* `serviceName` - A name for the service that you are creating.
* `environmentName` - The name of an existing Layer0 environment.
* `deployName[:deployVersion]` - The name of a Layer0 deploy that exists in the environment `environmentName`. You can optionally specify the version number of the Layer0 deploy that you want to deploy. If you do not specify a version number, the latest version of the deploy will be used.

####Optional arguments
* `--loadbalancer [environmentName:]loadBalancerName` - Place the new service behind an existing load balancer `loadBalancerName`. You can optionally specify the Layer0 environment (`environmentName`) where the Load Balancer exists.
* `--no-logs` - Disable cloudwatch logging for the service

### service update
Use the `update` subcommand to apply an existing Layer0 Deploy to an existing Layer0 service.

#### Usage
```
l0 service update [--no-logs] [environmentName:]serviceName deployName[:deployVersion]
```

####Required parameters
* `[environmentName:]serviceName` - The name of an existing Layer0 service into which you want to apply the deploy. You can optionally specify the Layer0 environment (`environmentName`) of the service.
* `deployName[:deployVersion]` - The name of the Layer0 deploy that you want to apply to the service. You can optionally specify a specific version of the deploy (`deployVersion`). If you do not specify a version number, the latest version of the deploy will be applied.

####Optional arguments
* `--no-logs` - Disable cloudwatch logging for the service

####Additional information
If your service uses a load balancer, when you update the task definition for the service, the container name and container port that were specified when the service was created must remain the same in the task definition. In other words, if your service has a load balancer, you cannot apply any deploy you want to that service. If you are varying the container name or exposed ports, you must create a new service instead.

###service delete
Use the `delete` subcommand to delete an existing Layer0 service.

#### Usage
```
l0 service delete [--wait] [environmentName:]serviceName
```

####Required parameters
* `[environmentName:]serviceName` - The name of the Layer0 service that you want to delete. You can optionally provide the Layer0 environment (`environmentName`) of the service.

####Optional arguments
* `--wait` - Wait until the deletion is complete before exiting.

####Additional information
This operation performs several tasks asynchronously. When run without the `--wait` option, this operation will most likely exit before all of these tasks are complete; when run with the `--wait` option, this operation will only exit once these tasks have completed.

###service get
Use the `get` subcommand to display information about an existing Layer0 service.

####Usage
```
l0 service get [environmentName:]serviceName
```

####Required parameters
* `[environmentName:]serviceName` - The name of an existing Layer0 service. You can optionally provide the Layer0 environment (`environmentName`) of the service.

###service list
Use the `list` subcommand to list all of the existing services in your Layer0 instance.

####Usage
```
l0 service get list
```

### service logs
Use the `logs` subcommand to display the logs from a Layer0 service that is currently running.

####Usage
```
l0 service logs [--start MM/DD HH:MM | --end MM/DD HH:MM | --tail=N] serviceName
```

####Required parameters
* `serviceName` - The name of the Layer0 service for which you want to view logs.


####Optional arguments
* `--start MM/DD HH:MM` - The start of the time range to fetch logs.
* `--end MM/DD HH:MM` - The end of the time range to fetch logs.
* `--tail=N` - Display only the last `N` lines of the log.

### service scale
Use the `scale` subcommand to specify how many copies of an existing Layer0 service should run.

####Usage
```
l0 service scale [environmentName:]serviceName copies
```

####Required parameters
* `[environmentName:]serviceName` - The name of the Layer0 service that you want to scale up. You can optionally provide the Layer0 environment (`environmentName`) of the service.
* `copies` - The number of copies of the specified service that should be run.

---

## Task
A Layer0 task is a component of an environment. A task executes the contents of a Docker image, as specified in a deploy. A task differs from a service in that a task does not restart after exiting. Additionally, ports are not exposed when using a task.

The `task` command is used with the following subcommands: [create](#task-create), [delete](#task-delete), [get](#task-get), [list](#task-list), and [logs](#task-logs).

### task create
Use the `create` subcommand to create a Layer0 task.

#### Usage
```
l0 task create [--copies copies | --no-logs] environmentName taskName deployName
```

####Required parameters
* `environmentName` - The name of the existing Layer0 environment in which you want to create the task.
* `taskName` - A name for the task.
* `deployName` - The name of an existing Layer0 deploy that the task should use.

####Optional arguments
* `--copies copies` - The number of copies of the task to run (default: 1).
* `--no-logs` - Disable cloudwatch logging for the service.

### task delete
Use the `delete` subcommand to delete an existing Layer0 task.

#### Usage
```
l0 task delete [environmentName:]taskName
```

####Required parameters
* `[environmentName:]taskName` - The name of the Layer0 task that you want to delete. You can optionally specify the name of the Layer0 environment that contains the task. This parameter is only required if mulitiple environments contain tasks with exactly the same name.

#### Additional information
Until the record has been purged, the API may indicate that the task is still running. Task records are typically purged within an hour.

### task get
Use the `get` subcommand to display information about an existing Layer0 task (`taskName`).

#### Usage
```
l0 task get [environmentName:]taskName
```

####Required parameters
* `[environmentName:]taskName` - The name of a Layer0 task for which you want to see information. You can optionally specify the name of the Layer0 Environment that contains the task.

####Additional information
The value of `taskName` does not need to exactly match the name of an existing task. If multiple results are found that match the pattern you specified in `taskName`, then information about all matching tasks will be returned.

### task list
Use the `task` subcommand to display a list of running tasks in your Layer0.

#### Usage
```
l0 task list
```

### task logs
Use the `logs` subcommand to display logs for a running Layer0 task.

#### Usage
```
l0 task logs [--start MM/DD HH:MM | --end MM/DD HH:MM | --tail=N] taskName
```

####Required parameters
* `taskName` - The name of an existing Layer0 task.

####Optional arguments
* `--start MM/DD HH:MM` - The start of the time range to fetch logs.
* `--end MM/DD HH:MM` - The end of the time range to fetch logs.
* `--tail=N` - Display only the last `N` lines of the log.

####Additional information
The value of `taskName` does not need to exactly match the name of an existing task. If multiple results are found that match the pattern you specified in `taskName`, then information about all matching tasks will be returned.

### task list
Use the `list` subcommand to display a list of running tasks in your Layer0.

#### Usage
```
l0 task list
```
