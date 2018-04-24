package aws

const DefaultLinuxUserdataTemplate = `
#!/bin/bash
echo ECS_CLUSTER={{ .ECSEnvironmentID }} >> /etc/ecs/ecs.config
echo ECS_ENGINE_AUTH_TYPE=dockercfg >> /etc/ecs/ecs.config
echo ECS_LOGLEVEL=info >> /etc/ecs/ecs.config
yum install -y aws-cli awslogs jq
aws s3 cp s3://{{ .S3Bucket }}/bootstrap/dockercfg dockercfg
cfg=$(cat dockercfg)
echo ECS_ENGINE_AUTH_DATA=$cfg >> /etc/ecs/ecs.config
docker pull amazon/amazon-ecs-agent:latest
start ecs
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
