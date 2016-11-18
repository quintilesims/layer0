# l0 command-line interface reference

##Global options

The **l0** application is designed to be used with one of several subcommands; these subcommands are detailed in the sections below. There are, however, some global parameters that you may specify when using **l0**.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0** [_globalOptions_] _command_ _subcommand_ [_options_] [_parameters_]</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--output {text|json}</div>
    <div class="divCell">Specify the format of Layer0 outputs. By default, Layer0 outputs unformatted text; by issuing the **--output json** option, you can force **l0** to output JSON-formatted text.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--version</div>
    <div class="divCell">Display the version number of the **l0** application.</div>
  </div>
</div>

---

##Admin
The **admin** command is used to manage the Layer0 API server. This command is used with the following subcommands: [debug](#admin-debug), [sql](#admin-sql), and [version](#admin-version).

### admin debug
Use the **debug** subcommand to view the running version of your Layer0 API server and CLI.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 admin debug**</div>
  </div>
</div>

### admin sql
Use the **sql** subcommand to initialize the Layer0 API database.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 admin sql**</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The **sql** subcommand is automatically executed during the Layer0 installation process; we recommend that you do not use this subcommand unless specifically directed to do so.</div>
  </div>
</div>

### admin version
Use the **version** subcommand to display the current version of the Layer0 API.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 admin version**</div>
  </div>
</div>

---

##Certificate
In order to use HTTPS ports in a Layer0 load balancer, you must create a certificate. You can use the Layer0 **certificate** command to upload and manage these certificates. This command is used with the following subcommands: [create](#certificate-create), [delete](#certificate-delete), [get](#certificate-get) and [list](#certificate-list).

### certificate create
Use the **create** subcommand to upload a certificate into Layer0. Once you have uploaded the certificate, you can use it when configuring HTTPS ports on Layer0 load balancers.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 certificate create** _certificateName publicKeyPath privateKeyPath_ [_intermediateChainPath_ ]</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_certificateName_</div>
    <div class="divCell">A name for the certificate.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_publicKeyPath_</div>
    <div class="divCell">The path to a public key associated with the certificate.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_privateKeyPath_</div>
    <div class="divCell">The path to the private key associated with the public key.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">*intermediateChainPath*</div>
    <div class="divCell">The path to an intermediate certificate authority.</div>
  </div>
</div>

### certificate delete
Use the **delete** subcommand to delete a certificate that has already been uploaded to Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 certificate delete** _certificateName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_certificateName_</div>
    <div class="divCell">The name of the certificate you want to delete.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">You will not be able to use the **delete** subcommand if a load balancer that uses that certificate is currently running.</div>
  </div>
</div>
  
### certificate get
Use the **get** subcommand to display information about an existing certificate in Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 certificate get** _certificateName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_certificateName_</div>
    <div class="divCell">The name of the Layer0 certificate for which you want to view additional information.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _certificateName_ does not need to exactly match the name of an existing certificate. If multiple results are found that match the pattern you specified in _certificateName_, then information about all matching certificates will be returned.</div>
  </div>
</div>

### certificate list
Use the **list** subcommand to list all of the certificates used in an instance of Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 certificate list**</div>
  </div>
</div>

---

##Deploy

### deploy create
Use the **create** subcommand to upload a Docker task definition into Layer0. This command is used with the following subcommands: [create](#deploy-create), [apply](#deploy-apply),  [delete](#deploy-delete), [get](#deploy-get) and [list](#deploy-list).

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 deploy create** _dockerPath_ _deployName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_dockerPath_</div>
    <div class="divCell">The path to the Docker task definition that you want to upload.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_deployName_</div>
    <div class="divCell">A name for the deploy.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">If _deployName_ exactly matches the name of an existing Layer0 deploy, then the version number of that deploy will increase by 1, and the task definition you specified will replace the task definition specified in the previous version.</div>
  </div> <br />
  <div class="divRow">
    <div class="divCellNoPadding">If you use Visual Studio to modify or create your Dockerrun file, you may see an "Invalid Dockerrun.aws.json" error. This error is caused by the default encoding used by Visual Studio. See the ["Common issues" page](http://localhost:8000/troubleshooting/commonissues/#invalid-dockerrunawsjson-error-when-creating-a-deploy) for steps to resolve this issue.</div>
  </div>

</div>

### deploy delete
Use the **delete** subcommand to delete a version of a Layer0 deploy.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 deploy delete** _deployID_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_deployID_</div>
    <div class="divCell">The unique identifier of the version of the deploy that you want to delete. You can obtain a list of deployIDs for a given deploy by executing the following command: <span class="noBreak">**l0 deploy get** _deployName_</span></div>
  </div>
</div>

### deploy get
Use the **get** subcommand to view information about an existing Layer0 deploy.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 deploy get** _deployName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_deployName_</div>
    <div class="divCell">The name of the Layer0 deploy for which you want to view additional information.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _deployName_ does not need to exactly match the name of an existing deploy. If multiple results are found that match the pattern you specified in _deployName_, then information about all matching deploys will be returned.</div>
  </div>
</div>

### deploy list
Use the **list** subcommand to view a list of deploys in your instance of Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 deploy list**</div>
  </div>
</div>

---

## Environment
Layer0 environments allow you to isolate services and load balancers for specific applications.
The **environment** command is used to manage Layer0 environments. This command is used with the following subcommands: [create](#environment-create), [delete](#environment-delete), [get](#environment-get), [list](#environment-list), and [setmincount](#environment-setmincount).

### environment create
Use the **create** subcommand to create an additional Layer0 environment (_environmentName_).

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 environment create** [--size] [--min-count] [--user-data] _environmentName_ </div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">A name for the environment.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--size</div>
    <div class="divCell">The size of the EC2 instances to create in your environment (default: m3.medium).</div>
  </div>
    <div class="divRow">
    <div class="divCellNoWrap">--min-count</div>
    <div class="divCell">The minimum number of EC2 instances allowed in the environment's autoscaling group (default: 0).</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--user-data</div>
    <div class="divCell">The user data template to use for the environment's autoscaling group.</div>
  </div>
</div>

The user data template can be used to add custom configuration to your Layer0 environment. 
Layer0 uses [Go Templates](https://golang.org/pkg/text/template) to render user data. 
Currently, two variables are passed into the template: **ECSEnvironmentID** and **S3Bucket**.
Please review the [ECS Tutorial](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_container_instance.html)
to better understand how to write a user data template, and use at your own risk! 
The default Layer0 user data template is:

```
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

### environment delete
Use the **delete** subcommand to delete an existing Layer0 environment.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 environment delete** [--wait] _environmentName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment that you want to delete.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--wait</div>
    <div class="divCell">Wait until the deletion is complete before exiting.</div>
  </div>
</div>

####Additional information 
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">This operation performs several tasks asynchronously. When run without the _--wait_ option, this operation will most likely exit before all of these tasks are complete; when run with the _--wait_ option, this operation will only exit once these tasks have completed.</div>
  </div>
</div>

### environment get
Use the **get** subcommand to display information about an existing Layer0 environment.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 environment get** _environmentName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment for which you want to view additional information.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _environmentName_ does not need to exactly match the name of an existing environment. If multiple results are found that match the pattern you specified in _environmentName_, then information about all matching environments will be returned.</div>
  </div>
</div>

### environment list
Use the **list** subcommand to display a list of environments in your instance of Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 environment list**</div>
  </div>
</div>

### environment setmincount
Use the **setmincount** subcommand to set the minimum number of EC2 instances allowed the environment's autoscaling group.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 enviroment setmincount** _environmentName_ _count_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment that you want to delete.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_count_</div>
    <div class="divCell">The minimum number of instances allowed in the environment's autoscaling group.</div>
  </div>
</div>

---

##Job
A Job is a long-running unit of work performed on behalf of the Layer0 API.
Jobs are executed as Layer0 tasks that run in the **api** Environment. 
The **job** command is used with the following subcommands: [logs](#job-logs), [delete](#job-delete), [get](#job-get), and [list](#job-list).

### job logs
Use the **logs** subcommand to display the logs from a Layer0 job that is currently running.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 job logs** [--tail=*N* ] _jobName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_jobName_</div>
    <div class="divCell">The name of the Layer0 job for which you want to view logs.</div>
  </div>
</div>
### job logs
Use the **logs** subcommand to display the logs from a Layer0 job that is currently running.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 job logs** [--tail=*N* ] _jobName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_jobName_</div>
    <div class="divCell">The name of the Layer0 job for which you want to view logs.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--tail=*N*</div>
    <div class="divCell">Display only the last _N_ lines of the log.</div>
  </div>
</div>

###job delete
Use the **delete** subcommand to delete an existing job.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 job delete** *jobName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">*jobName*</div>
    <div class="divCell">The name of the job that you want to delete.</div>
  </div>
</div>

###job get
Use the **get** subcommand to display information about an existing Layer0 job.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 job get** *jobName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_jobName_</div>
    <div class="divCell">The name of an existing Layer0 job.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _jobName_ does not need to exactly match the name of an existing job. If multiple results are found that match the pattern you specified in _jobName_, then information about all matching jobs will be returned.</div>
  </div>
</div>

###job list
Use the **list** subcommand to display information about all of the existing jobs in an instance of Layer0.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 job list**</div>
  </div>
</div>

---

##Loadbalancer
A load balancer is a component of a Layer0 environment. Load balancers listen for traffic on certain ports, and then forward that traffic to Layer0 [services](#service). The **loadbalancer** command is used with the following subcommands: [create](#loadbalancer-create), [delete](#loadbalancer-delete), [addport](#loadbalancer-addport), [dropport](#loadbalancer-dropport), [get](#loadbalancer-get), and [list](#loadbalancer-list).

###loadbalancer create
Use the **create** subcommand to create a new load balancer.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer create** [--port _port_ --port _port_ ...] [--certificate _certificateName_] [--private] _environmentName loadBalancerName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the existing Layer0 environment in which you want to create the load balancer.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_loadBalancerName_</div>
    <div class="divCell">A name for the load balancer.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--port _hostPort:containerPort/protocol_</div>
    <div class="divCell">The port configuration for the load balancer. _hostPort_ is the port on which the load balancer will listen for traffic; _containerPort_ is the port that traffic will be forwarded to. You can specify multiple ports using _--port xxx --port yyy_. If this option is not specified, Layer0 will use the following configuration: 80:80/tcp</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--certificate _certificateName_</div>
    <div class="divCell">The name of an existing Layer0 certificate. You must include this option if you are using an HTTPS port configuration.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--private</div>
    <div class="divCell">When you use this option, the load balancer will only be accessible from within the Layer0 environment.</div>
  </div>
</div>

###loadbalancer delete
Use the **delete** subcommand to delete an existing load balancer.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer delete** [--wait] *loadBalancerName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">*loadBalancerName*</div>
    <div class="divCell">The name of the load balancer that you want to delete.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--wait</div>
    <div class="divCell">Wait until the deletion is complete before exiting.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">In order to delete a load balancer that is already attached to a service, you must first delete the service that uses the load balancer.</div>
  </div><br />
  <div class="divRow">
    <div class="divCellNoPadding">This operation performs several tasks asynchronously. When run without the _--wait_ option, this operation will most likely exit before all of these tasks are complete; when run with the _--wait_ option, this operation will only exit once these tasks have completed.</div>
  </div>
</div>

###loadbalancer addport
Use the **addport** subcommand to add a new port configuration to an existing Layer0 load balancer.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer addport** *loadBalancerName hostPort:containerPort/protocol* [--certificate _certificateName_]</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_loadBalancerName_</div>
    <div class="divCell">The name of an existing Layer0 load balancer in which you want to add the port configuration.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_hostPort_</div>
    <div class="divCell">The port that the load balancer will listen on.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_containerPort_</div>
    <div class="divCell">The port that the load balancer will forward traffic to.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_protocol_</div>
    <div class="divCell">The protocol to use when forwarding traffic (acceptable values: tcp, ssl, http, and https).</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--certificate _certificateName_</div>
    <div class="divCell">The name of an existing Layer0 certificate. You must include this option if you are using an HTTPS port configuration.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The port configuration you specify must not already be in use by the load balancer you specify.</div>
  </div>
</div>

###loadbalancer dropport
Use the **dropport** subcommand to remove a port configuration from an existing Layer0 load balancer. 

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer dropport** *loadBalancerName* *hostPort*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_loadBalancerName_</div>
    <div class="divCell">The name of an existing Layer0 load balancer in which you want to remove the port configuration.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_hostPort_</div>
    <div class="divCell">The host port to remove from the load balancer.</div>
  </div>
</div>

###loadbalancer get
Use the **get** subcommand to display information about an existing Layer0 load balancer. 

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer get** *environmentName:loadBalancerName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of an existing Layer0 environment.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_loadBalancerName_</div>
    <div class="divCell">The name of an existing Layer0 load balancer.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _loadBalancerName_ does not need to exactly match the name of an existing load balancer. If multiple results are found that match the pattern you specified in _loadBalancerName_, then information about all matching load balancers will be returned.</div>
  </div>
</div>

###loadbalancer list
Use the **list** subcommand to display information about all of the existing load balancers in an instance of Layer0. 

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 loadbalancer list**</div>
  </div>
</div>

---

## Service
A service is a component of a Layer0 environment. The purpose of a service is to execute a Docker image specified in a [deploy](#deploy). In order to create a service, you must first create an [environment](#environment) and a [deploy](#deploy); in most cases, you should also create a [load balancer](#loadbalancer) before creating the service.

The **service** command is used with the following subcommands: [create](#service-create), [delete](#service-delete), [get](#service-get), [list](#service-list), [logs](#service-logs), and [scale](#service-scale).

###service create
Use the **create** subcommand to create a Layer0 service. 

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service create** [--loadbalancer _environmentName:loadBalancerName_ ] [--no-logs] _environmentName serviceName deployName:deployVersion_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">A name for the service that you are creating.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of an existing Layer0 environment.</div>
  </div>  
  <div class="divRow">
    <div class="divCellNoWrap">_deployName_</div>
    <div class="divCell">The name of a Layer0 deploy that exists in the environment _environmentName_.</div>
  </div> 
  <div class="divRow">
    <div class="divCellNoWrap">_deployVersion_</div>
    <div class="divCell">The version number of the Layer0 deploy that you want to deploy. If you do not specify a version number, the latest version of the deploy will be used.</div>
  </div> 
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--loadbalancer _environmentName:loadBalancerName_</div>
    <div class="divCell">Place the new service behind an existing load balancer named _loadBalancerName_ in the environment named _environmentName_.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--no-logs</div>
    <div class="divCell">Disable cloudwatch logging for the service</div>
  </div>
</div>

### service update
Use the **update** subcommand to apply an existing Layer0 Deploy to an existing Layer0 service.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service update** [--no-logs] _environmentName:serviceName deployName:deployVersion_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
 <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment in which the service resides.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">The name of an existing Layer0 service into which you want to apply the deploy.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_deployName_</div>
    <div class="divCell">The name of the Layer0 deploy that you want to apply to the service.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_deployVersion_</div>
    <div class="divCell">The version of the Layer0 deploy that you want to apply to the service. If you do not specify a version number, the latest version of the deploy will be applied.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--no-logs</div>
    <div class="divCell">Disable cloudwatch logging for the service</div>
  </div>
</div>

####Additional information

<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">If your service uses a load balancer, when you update the task definition for the service, the container name and container port that were specified when the service was created must remain the same in the task definition. In other words, if your service has a load balancer, you cannot apply any deploy you want to that service. If you are varying the container name or exposed ports, you must create a new service instead.</div>
  </div>
</div>


###service delete
Use the **delete** subcommand to delete an existing Layer0 service.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service delete** [--wait] _environmentName:serviceName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment that contains the service you want to delete.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">The name of the Layer0 service that you want to delete.</div>
  </div> 
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--wait</div>
    <div class="divCell">Wait until the deletion is complete before exiting.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">This operation performs several tasks asynchronously. When run without the _--wait_ option, this operation will most likely exit before all of these tasks are complete; when run with the _--wait_ option, this operation will only exit once these tasks have completed.</div>
  </div>
</div>

###service get
Use the **get** subcommand to display information about an existing Layer0 service.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service get** _environmentName:serviceName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of an existing Layer0 environment.</div>
  </div> 
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">The name of an existing Layer0 service.</div>
  </div> 
</div>

###service list
Use the **list** subcommand to list all of the existing services in your Layer0 instance.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service list**</div>
  </div>
</div>

### service logs
Use the **logs** subcommand to display the logs from a Layer0 service that is currently running.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service logs** [--tail=*N* ] _serviceName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">The name of the Layer0 service for which you want to view logs.</div>
  </div>
</div>

####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--tail=*N*</div>
    <div class="divCell">Display only the last _N_ lines of the log.</div>
  </div>
</div>

### service scale
Use the **scale** subcommand to specify how many copies of an existing Layer0 service should run.

####Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 service scale** _environmentName:serviceName N_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the Layer0 environment that contains the service that you want to scale.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_serviceName_</div>
    <div class="divCell">The name of the Layer0 service that you want to scale up.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_N_</div>
    <div class="divCell">The number of copies of the specified service that should be run.</div>
  </div>
</div>

---

## Task
A Layer0 task is a component of an environment. A task executes the contents of a Docker image, as specified in a deploy. A task differs from a service in that a task does not restart after exiting. Additionally, ports are not exposed when using a task.

The **task** command is used with the following subcommands: [create](#task-create), [delete](#task-delete), [get](#task-get), [list](#task-list), and [logs](#task-logs).

### task create
Use the **create** subcommand to create a Layer0 task.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task create** [--no-logs] [--copies _copies_] *environmentName taskName deployName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_environmentName_</div>
    <div class="divCell">The name of the existing Layer0 environment in which you want to create the task.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_taskName_</div>
    <div class="divCell">A name for the task.</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">_deployName_</div>
    <div class="divCell">The name of an existing Layer0 deploy that the task should use.</div>
  </div>
</div>
  
####Optional arguments
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">--copies</div>
    <div class="divCell">The number of copies of the task to run (default: 1)</div>
  </div>
  <div class="divRow">
    <div class="divCellNoWrap">--no-logs</div>
    <div class="divCell">Disable cloudwatch logging for the service</div>
  </div>
</div>

### task delete
Use the **delete** subcommand to delete an existing Layer0 task.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task delete** [*environmentName*:]*taskName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_taskName_</div>
    <div class="divCell">The name of the Layer0 task that you want to delete.</div>
  </div>
</div>

####Optional parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">[_environmentName_:]</div>
    <div class="divCell">The name of the Layer0 environment that contains the task. This parameter is only necessary if multiple environments contain tasks with exactly the same name.</div>
  </div>
</div>

#### Additional information 
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">Until the record has been purged, the API may indicate that the task is still running. Task records are typically purged within an hour.</div>
  </div>
</div>

### task get
Use the **get** subcommand to display information about an existing Layer0 task (_taskName_).

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task get** [*environmentName*:]*taskName*</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_taskName_</div>
    <div class="divCell">The name of a Layer0 task for which you want to see information.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _taskName_ does not need to exactly match the name of an existing task. If multiple results are found that match the pattern you specified in _taskName_, then information about all matching tasks will be returned.</div>
  </div>
</div>

### task list
Use the **task** subcommand to display a list of running tasks in your Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task list**</div>
  </div>
</div>

### task logs
Use the **logs** subcommand to display logs for a running Layer0 task.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task logs** [--tail=_N_ ] _taskName_</div>
  </div>
</div>

####Required parameters
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoWrap">_taskName_</div>
    <div class="divCell">The name of an existing Layer0 task.</div>
  </div>
</div>

####Additional information
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">The value of _taskName_ does not need to exactly match the name of an existing task. If multiple results are found that match the pattern you specified in _taskName_, then information about all matching tasks will be returned.</div>
  </div>
</div>

### task list
Use the **task** subcommand to display a list of running tasks in your Layer0.

#### Usage
<div class="divTable">
  <div class="divRow">
    <div class="divCellNoPadding">**l0 task list**</div>
  </div>
</div>
=======
