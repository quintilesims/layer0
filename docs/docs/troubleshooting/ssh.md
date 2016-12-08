# Secure Shell (SSH)
You can use Secure Shell (SSH) to access your Layer0 environment(s).

By default, Layer0 Setup asks for an EC2 key pair when creating a new Layer0. This key pair is associated with all machines that host your Layer0 Services. This means you can use SSH to log into the underlying Docker host to perform tasks such as troubleshooting failing containers or viewing logs. For information about creating an EC2 key pair, see [Install and Configure Layer0](/setup/install.md#Prerequisites).

!!! warning
	This section is recommended for development debugging only.
	It is **not** recommended for production environments.

# To SSH into a Service
1. In a console window, add port 2222:22/tcp to your Service's load balancer:
```
l0 loadbalancer addport <name> 2222:22/tcp
```
<ol start="2">
  <li>SSH into your Service by supplying the load balancer url and key pair file name.</li>
</ol>
```
ssh -i <key pair path and file name> ec2-user@<load balancer url> -p 2222
```
<ol start="3">
  <li>If required, Use Docker to access a specific container with Bash.</li>
</ol>
```
docker exec -it <container id> /bin/bash
```
##Remarks
You can get the load balancer url from the Load Balancers section of your Layer0 AWS console.

Use the [loadbalancer dropport](/reference/cli.md#loadbalancer) subcommand to remove a port configuration from an existing Layer0 load balancer.

You _cannot_ change the key pair after a Layer0 has been created. If you lose your key pair or need to generate a new one, you will need to create a new Layer0.

If your Service is behind a private load balancer, or none at all, you can either re-create your Service behind a public load balancer, use an existing public load balancer as a "jump" point, or create a new Layer0 Service behind a public load balancer to serve as a "jump" point.
