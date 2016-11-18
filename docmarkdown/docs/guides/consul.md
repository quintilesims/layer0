# Deployment guide: Consul service

Consul is a tool for configuring services in your infrastructure. It includes several features, including health monitoring, service discovery and key/value storage.

This guide provides step-by-step instructions for deploying Consul in a Layer0 instance. These procedures build upon the [Guestbook](/guides/guestbook) and [Guestbook with RDS](/guides/guestbook_rds) deployment guides; you must complete the procedures in those guides before you can complete the procedures in this guide.

## Before you start
In order to complete the procedures in this section, you must install and configure Layer0 v0.7.0 or later. If you have not already configured Layer0, see the [installation guide](/setup/install). If you are running an older version of Layer0, see the [upgrade instructions](/setup/update#upgrading-older-versions-of-layer0).

This guide expands upon the [Guestbook with RDS Database deployment guide](/guides/guestbook_rds). The procedures in this guide assume that you completed the Guestbook with RDS deployment guide and all of its prerequisites.

## Part 1: Create a load balancer

Consul should run behind a private load balancer in the **demo** environment with ports 8500 and 8301 exposed.

**To create the load balancer:**
<ol>
  <li>At the command line, type the following command to create a load balancer named **consullb** in the **demo** environment with port 8500 exposed:
    <ul>
      <li class="command">**l0 loadbalancer create --private --port 8500:8500/tcp demo consullb**</li>
    </ul><br />You will see the following output:
<pre class="code"><code>LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICES  PORTS          PUBLIC  URL
1consullb        consullb           demo                   8500:8500/tcp  false   
</code></pre>
  </li>
  <li>At the command line, type the following command to add port 8301 to the **consullb** load balancer:
    <ul>
      <li class="command">**l0 loadbalancer addport demo:consullb 8301:8301/tcp**</li>
    </ul><br />
  You will see the following output:
<pre class="code"><code>LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICES  PORTS          PUBLIC  URL
1consullb        consullb           demo                   8500:8500/tcp  false   <em>(url)</em>
                                                           8301:8301/tcp
</code></pre>
<br />
Copy the URL listed in the **URL** column; you will need this URL in the next section.</li>
</ol>

## Part 2: Configure the deploy
<ol>
  <li>Download the <a href="https://gitlab.imshealth.com/xfra/consul-local/blob/master/consul.json">Consul Task Definition</a> and save it to your computer as Consul.Dockerrun.aws.json.</li>
  <li>Open Consul.Dockerrun.aws.json in a text editor. Toward the bottom of the file, you will see the following:
<pre class="code"><code>"environment": [
    {
        "name": "EXTERNAL_URL",
        "value": "&lt;url&gt;"
     }
]
</code></pre>
<br />
Replace <em>&lt;url&gt;</em> with the URL you copied in step 2 of the previous section, and then save the file.</li>
  <li>At the command line, type the following command to create a new deploy called <strong>consul</strong>:
    <ul>
      <li class="command"><strong>l0 deploy create Consul.Dockerrun.aws.json consul</strong></li>
    </ul><br />
  You will see the following output:
<pre class="code"><code>DEPLOY ID  DEPLOY NAME  VERSION  
1consul:1  consul       1        
</code></pre>
  </li>
</ol>

## Part 3: Create the service
Now that you've created an environment, load balancer, and deploy, you can create a service to bring these elements together.

**To create the service:**
<ol>
  <li>At the command line, type the following command to create a new service called **consul**:
    <ul>
      <li class="command">**l0 service create --loadbalancer demo:consullb demo consulsvc consul:latest**</li>
    </ul><br />
    You will see the following output:
<pre class="code"><code>SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYS      SCALE
1consul     consulsvc     demo         consullb      consul:1     0/1</code></pre>
  </li>
  <li>Wait several minutes for the service to be provisioned. You can check the status of the service creation by running the following command:
  <ul>
    <li class="command">**l0 service get consulsvc**</li>
  </ul><br />
  When the service has finished provisioning, you will see the following output:
```
SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYS      SCALE
1consul     consulsvc     demo         consullb      consul:1     1/1
```
  </li>
</ol>

## Part 4: Scale the Consul service
Consul is a scalable application. For added reliability, we recommend that you scale the Consul service to size 3.

**To scale the service:**

<ol>
  <li>At the command line, type the following to scale the consul service to size 3:
    <ul>
      <li class="command"><strong>l0 service scale demo:consulsvc 3</strong></li>
    </ul>
  </li>

  <li>Wait several minutes for the service to scale. You can check the status of the service by running the following command:
  <ul>
    <li class="command"><strong>l0 service get consul</strong></li>
  </ul><br />
  When the Service has finished scaling, you will see the following output:
<pre class="code"><code>SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYS      SCALE
1consul     consulsvc     demo         consullb      consul:1     3/3
</code></pre>
  </li>
</ol>

## Additional steps: Configure Layer0 services
In order to for Layer0 services to use consul, the task definitions for those services must be configured. The next deployment guide in this series ([Guestbook with Consul](/guides/guestbook_consul)) contains an example of these configurations.
