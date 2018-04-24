Layer0 Terraform Module
===========
A Terraform module to provide a Layer0 instance in AWS.


Module Input Variables
----------------------

- `name` - Layer0 instance name
- `layer0_version` - Version of Layer0 to use
- `access_key` - AWS Access Key ID to manage resources
- `secret_key` - AWS Secret Access Key to manage resources
- `region` - AWS Region to manage resources
- `ssh_key_pair` - Name of an existing AWS SSH key pair used to SSH into EC2 instances
- `dockercfg` - Docker authentication contents in `dockercfg` format
- `username` - Username to use when accessing the API
- `password` - Password to use when accessing the API
- `vpc_id` - (Optional) ID of an existing VPC to install the Layer0 instance.
If none specified, a new VPC will be created


Usage
-----

```hcl
module "layer0" {
  source         = "https://github.com/quintilesims/layer0/setup//layer0?ref=v1.0.0"
  name           = "foobar"
  layer0_version = "v1.0.0"
  access_key     = "ABC123"
  secret_key     = "ABC123"
  region         = "us-west-2"
  ssh_key_pair   = "my-key"
  dockercfg      = "${file("~/.docker/dockercfg")}"
  username       = "layer0"
  password       = "password123"
  vpc_id         = ""
}
```


Outputs
=======
- `name` - Layer0 instance name
- `account_id` - AWS Account ID in which the Layer0 instance was created
- `endpoint` - Endpoint of the API
- `token` - Authentication token to use when communicating with the API
- `s3_bucket` - Name of the S3 bucket created for the Layer0 instance
- `access_key` - AWS Access Key ID generated for the API 
- `secret_key` - AWS Secret Access Key generated for the API 
- `vpc_id` - ID of the VPC in which the Layer0 instance was created
- `public_subnets` - Comma-separated list of public subnet IDs in the VPC
- `private_subnets` - Comma-separated list of private subnet IDs in the VPC
- `ecs_role` - Name of the IAM role created for the Layer0 Instance
- `ssh_key_pair` - Name of the SSH key pair used by the API
- `ecs_agent_instance_profile` - Name of the IAM instance profile created for the Layer0 instance
- `linux_service_ami` - The AMI used when creating Linux EC2 instances

