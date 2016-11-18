# Deployment guide: Guestbook with Consul
This guide provides step-by-step instructions for deploying a Guestbook application that stores data in a
[Redis](http://redis.io) database. The Guestbook application uses [Consul](/reference/consul) to dynamically discover the Redis service.

## Before you start
This guide assumes that you are running Layer0 version 0.7.2 or later, and that you have completed the
[Guestbook](/guides/guestbook) and [Consul](/guides/consul) deployment guides.

## Part 1: Configure and deploy the Redis task definition
The updated Guestbook application in this guide stores its data in a Redis database. Before you can deploy the updated Guestbook application, you must first configure and deploy the Redis task definition.

**To configure and deploy the task definition:**
<ol>
  <li>Download the [Redis task
definition](https://gitlab.imshealth.com/xfra/layer0-samples/blob/master/redis/Redis.Dockerrun.aws.json) and save it to your computer as Redis.dockerrun.aws.json.</li>
  <li>At the command prompt, type the following command:
    <ul>
      <li class="command"><strong>l0 loadbalancer get consul</strong></li>
    </ul><br />
    Copy the value in the <strong>URL</strong> column; you will need it in the next step.
  </li>
  <li>Open Redis.dockerrun.aws.json in a text editor. Toward the end of the file, you will see the following:
<pre class="code"><code>"environment": [
    {
        "name": "EXTERNAL_URL",
        "value": "&lt;url&gt;"
    }
]</code></pre>
    In this section, replace <em>&lt;url&gt;</em> with the URL you copied in the previous step.
  </li>
  <li>At the command prompt, type the following command to create a new deploy named **redis** that uses the Redis.dockerrun.aws task definition:
    <ul>
      <li class="command">**l0 deploy create Redis.Dockerrun.aws.json redis**</li>
    </ul><br />
  You will see the following output:
<pre class="code"><code>DEPLOY ID  DEPLOY NAME  VERSION  
1redis:1   redis        1</code></pre>  
  </li>  
</ol>

##Part 2: Create the redis service
Now that you have created the **redis** deploy, you can create a service to place it in.

**To create a new service:**
<ul>
  <li>At the command line, type the following command to create a service named **redis** running the deploy you created in the previous section:
    <ul>
      <li class="command"> **l0 service create demo redis redis:latest**</li>
    </ul><br />
    When you execute this command, the [Registrator](/reference/consul/#registrator) service will read the environment variables for the l0-demo-redis container. These variables include the following:

```
"environment": [
    {
        "name": "SERVICE_NAME",
        "value": "db"
    },
    {
        "name": "SERVICE_TAGS",
        "value": "guestbook"
    },
    ...
```
<br />After reading these variables, Registrator will register a service named **db** into your environment's Consul service.
Members of the Consul Cluster will be able to discover this service via DNS queries to guestbook.db.service.consul.
  </li>
</ul>

##Part 3: Update the Guestbook service

To configure the Guestbook application to use Consul for service discovery, you must first update the **guestbook** deploy.

**To update the Guestbook service:**

<ol>
  <li>Download the <a href="https://gitlab.imshealth.com/xfra/layer0-samples/blob/master/redis/Guestbook.Dockerrun.aws.json">Guestbook with Consul Task definition</a> and save it to your computer as GuestbookConsul.Dockerrun.aws.json.</li>
  <li>At the command line, type the following command: <strong>l0 loadbalancer get consullb</strong>. Copy the value in the <strong>URL</strong> column.</li>
  <li>Open GuestbookConsul.Dockerrun.aws.json in a text editor. Toward the bottom of the file, in the <strong>environment</strong> section, replace <em>&lt;url&gt;</em> with the URL
that you copied in the previous step. Save the file.</li>
  <li>At the command line, type the following command to create a new version of the guestbook deploy using the updated Guestbook task definition:
    <ul>
      <li class="command"><strong>l0 deploy create GuestbookConsul.Dockerrun.aws.json guestbook</strong></li>
    </ul><br />
  You will see the following output:
  
<pre class="code"><code>DEPLOY ID     DEPLOY NAME  VERSION
1guestbook:3  guestbook    3
</code></pre>
  </li>
  <li>At the command line, type the following command to apply the updated <strong>guestbook</strong> deploy on the <strong>guestbook</strong> service:
    <ul>
      <li class="command"><strong>l0 deploy apply guestbook:latest demo:guestbook</strong></li>
    </ul><br />
  When you execute this command, the Guestbook application will automatically discover the Redis service by making a DNS query to guestbook.db.service.consul.</li>
  <li>At the command prompt, type the following command to find the URL of the <strong>guestbooklb</strong> load balancer:
    <ul>
      <li class="command"><strong>l0 loadbalancer get guestbooklb</strong></li>
    </ul><br />
  If the service has not finished deploying, you may see the following message:
"Could not connect to redis server, try refreshing."<br />
If you see this message, wait a few minutes, and then refresh the page. You can also try using your web browser's Incognito or Private Browsing mode to ensure that the page is not cached.
  </li>
</ol>
