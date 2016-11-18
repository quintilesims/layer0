#!/bin/bash
echo ECS_CLUSTER={{ .ECSEnvironmentID }} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
yum install -y aws-cli awslogs jq
aws s3 cp s3://{{ .S3Bucket }}/bootstrap/dockercfg dockercfg
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs
