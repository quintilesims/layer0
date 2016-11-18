# Deployment guide: Shinken
Shinken is an open-source application monitoring framework that provides several benefits, including:

  * Externally-accessible status URLs
  * Health checks for applications that use Consul
  * Slack channel integration

This guide provides instructions for deploying a Shinken service using Layer0.

## Before You Start
This guide assumes you have an instance of Layer0 version v0.7.0 or higher; if not, complete the [Installation Instructions](/setup/install).

Once Layer0 is configured on your computer, download the [Shinken Task Definition](https://gitlab.imshealth.com/xfra/layer1-shinken/blob/master/Shinken.aws.json); name the resulting file **Shinken.Dockerrun.aws.json**.

This guide expands upon the [Guestbook with Consul](/guides/guestbook_consul) deployment guide. You must complete the procedures in that guide before you can complete the procedures listed here.

##Part 1: Create the load balancer
When run as a Layer0 service, Shinken should be placed behind a public load balancer. The load balancer should transfer traffic via port 80 over the HTTP protocol.

**To create the load balancer:**

<ol>
  <li>At the command prompt, type the following command to create a new load balancer named **shinkenlb** in the **demo** environment that forwards traffic through port 80 via the HTTP protocol:
    <ul>
      <li class="command"><strong>l0 loadbalancer create --port 80:80/http demo shinkenlb</strong></li>
    </ul><br />
You will see the following output:
<pre class="code"><code>LOADBALANCER ID  LOADBALANCER NAME  ENVIRONMENT  SERVICES  PORTS       PUBLIC  URL
<em>(id)</em>             shinkenlb          demo                   80:80/http  false   <em>(url)</em></code></pre>
  </li>
</ol>


## Part 2: Configure the deploy
Before you can deploy Shinken, you must modify the Shinken task definition.

**To configure the Shinken task definition:**
<ol>
  <li>In a text editor, open Shinken.Dockerrun.aws.json.</li>
  <li>Toward the bottom of the file, in the **environment** section, you will see the following:
  <pre class="code"><code>"environment": [
    {
        "name": "SLACK_WEBHOOK_URL",
        "value": ""
    },
    {
        "name": "EXTERNAL_URL",
        "value": ""
    },
    {
        "name": "ADMIN_PASSWORD_HASH",
        "value": ""
    }
]</code></pre><br />
  Modify the file to include the following values:
  <ul>
    <li><strong>SLACK_WEBHOOK_URL</strong>: The incoming webhook URL for your team's Slack channel.
      <ul>
        <li>To find the webhook URL for your channel, visit the <a href="https://ims-dev.slack.com/apps/manage/A0F7XDUAZ-incoming-webhooks">Incoming WebHooks</a> page, and then click the name of the channel. Copy the value shown in the **Webhook URL** field.</li>
        <li>Alternatively, if an incoming webhook does not already exist for your channel, click the **Add Configuration** button to create a new incoming webhook integration.</li>
        <li>The Xfra team has created a #test channel in Slack for testing webhooks. The webhook URL for the #test channel is https://hooks.slack.com/services/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</li>
      </ul>
    </li>
    <li><strong>EXTERNAL_URL</strong>: The URL for the <strong>shinkenlb</strong> load balancer. To find the URL, type the following command, and then copy the value in the <strong>URL</strong> column:
      <ul>
        <li class="command"><strong>l0 loadbalancer get shinkenlb</strong></li>
      </ul>
    </li>
    <li><strong>ADMIN_PASSWORD_HASH</strong>: A hashed password that you will use to access the Shinken user interface, in the format *admin:&lt;hash&gt;*
      <ul>
        <li>Linux and Mac users can type the following command in a terminal window to obtain a password hash:
        <ul>
          <li class="command"><strong>htpasswd -nb admin</strong> <em>yourPassword</em></li>
        </ul>
        <li>Windows users should use <a href="http://www.htaccesstools.com/htpasswd-generator/">an online htpasswd generator</a> to create the password hash.</li>
      </ul>
    </li><br />
    When you have finished making changes, save the file.
  </ul>
  <li>Type the following command to create a deploy named **shinken** using the task definition you just modified:
    <ul>
      <li class="command"><strong>l0 deploy create Shinken.Dockerrun.aws.json shinken</strong></li>
    </ul><br />
    You will see the following output:
<pre class="code"><code>DEPLOY ID   DEPLOY NAME  VERSION
1shinken:1  shinken      1</code></pre>
  </li>
</ol>

##Part 3: Create the service
Now that you have created the **shinken** deploy, you can use the **service** command to run it.

**To create the service:**
<ol>
  <li>At the command prompt, type the following command to create a service named **shinkensvc** that uses the **shinkenlb** load balancer:
    <ul>
      <li class="command"><strong>l0 service create --loadbalancer demo:shinkenlb demo shinkensvc shinken:latest</strong></li>
    </ul><br />
    You will see the following output:
<pre class="code"><code>SERVICE ID   SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYS      SCALE
1shinken     shinkensvc    demo         shinkenlb     shinken:1    0/1
</code></pre>
  </li>
  <li>Wait several minutes for the service to provision completely. Use the following command to check the status of the service provisioning:
    <ul>
      <li class="command"><strong>l0 service get demo:shinkensvc</strong></li>
    </ul><br />
    When the service is ready, you will see the following output:
<pre class="code"><code>SERVICE ID  SERVICE NAME  ENVIRONMENT  LOADBALANCER  DEPLOYS      SCALE
<em>(id)</em>        shinkensvc    demo         shinkenlb     shinken:1    1/1</code></pre>
  </li>
</ol>

##Additional steps: Configuring other services
In order to use Shinken health checks, you will need to modify the task definitions for your Layer0 applications. For more information, see the [Shinken Reference page](/reference/shinken/#service-configuration).
