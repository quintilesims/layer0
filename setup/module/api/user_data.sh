#!/bin/bash
echo ECS_CLUSTER=${cluster_id} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
yum install -y aws-cli awslogs
aws s3 cp s3://${s3_bucket}/bootstrap/dockercfg dockercfg
aws s3 sync s3://{s3_bucket}/terraform/ /usr/local/bin/ --exclude "*" --include "docker-credential-*"
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs
