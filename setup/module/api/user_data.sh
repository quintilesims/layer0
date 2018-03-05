# The majority of this script was largely taken from AWS's documentation
# on using CloudWatch logs with container instances
# https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_cloudwatch_logs.html

Content-Type: multipart/mixed; boundary="==BOUNDARY=="
MIME-Version: 1.0

--==BOUNDARY==
#!/bin/bash
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1
echo ECS_CLUSTER=${cluster_id} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
echo ECS_LOGLEVEL=info >> /etc/ecs/ecs.config

# Known issue where packages will sometimes not install on the first try
n=0
until [ $n -ge 5 ]
do
    yum install -y aws-cli awslogs jq && break
    n=$[$n+1]
    sleep 5
done

# Inject the CloudWatch Logs configuration file contents
cat > /etc/awslogs/awslogs.conf <<- EOF
[general]
state_file = /var/lib/awslogs/agent-state        
 
[/var/log/docker]
file = /var/log/docker
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/docker
datetime_format = %Y-%m-%dT%H:%M:%S.%f
[/var/log/messages]
file = /var/log/messages
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/messages
datetime_format = %Y-%m-%dT%H:%M:%S.%f
[/var/log/user-data.log]
file = /var/log/docker
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/user-data.log
datetime_format = %Y-%m-%dT%H:%M:%S.%f
[/var/log/ecs/ecs-init.log]
file = /var/log/ecs/ecs-init.log
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/ecs/ecs-init.log
datetime_format = %Y-%m-%dT%H:%M:%SZ
[/var/log/ecs/ecs-agent.log]
file = /var/log/ecs/ecs-agent.log.*
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/ecs/ecs-agent.log
datetime_format = %Y-%m-%dT%H:%M:%SZ
[/var/log/ecs/cloudwatch-logs-start.log]
file = /var/log/ecs/cloudwatch-logs-start.log
log_group_name = {cluster}
log_stream_name = {container_instance_id}/var/log/ecs/cloudwatch-logs-start.log
datetime_format = %Y-%m-%dT%H:%M:%SZ
EOF

aws s3 cp s3://${s3_bucket}/bootstrap/dockercfg dockercfg
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs

--==BOUNDARY==
Content-Type: text/x-shellscript; charset="us-ascii"
#!/bin/bash
# Set the region to send CloudWatch Logs data to (the region where the container instance is located)
region=$(curl -s 169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
sed -i -e "s/region = us-east-1/region = $region/g" /etc/awslogs/awscli.conf

--==BOUNDARY==
Content-Type: text/upstart-job; charset="us-ascii"

#upstart-job
description "Configure and start CloudWatch Logs agent on Amazon ECS container instance"
author "Amazon Web Services"
start on started ecs

script
	exec 2>>/var/log/ecs/cloudwatch-logs-start.log
	set -x

	until curl -s http://localhost:51678/v1/metadata
	do
		sleep 1	
	done

	# Grab container instance id from instance metadata
	cluster=$(curl -s http://localhost:51678/v1/metadata | jq -r '. | .Cluster' | cut -d '/' -f2 | cut -d '-' -f1,2)
	container_instance_id=$(curl -s 169.254.169.254/latest/dynamic/instance-identity/document | jq -r .instanceId)

	# Replace the container instance ID placeholder with the actual value
	sed -i -e "s#{cluster}#$cluster#g" /etc/awslogs/awslogs.conf
	sed -i -e "s#{container_instance_id}#$container_instance_id#g" /etc/awslogs/awslogs.conf

	service awslogs start
	chkconfig awslogs on
end script
--==BOUNDARY==--

