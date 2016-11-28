# Updating a Layer0 service
There are three methods of updating an existing Layer0 service. The first method is to update the existing Deploy to refer to a new Docker task definition. The second method is to create a new Service that uses the same Loadbalancer. The third method is to create both a new Loadbalancer and a new Service.

There are advantages and disadvantages to each of these methods. The following sections discuss the advantages and disadvantages of using each method, and include procedures for implementing each method.

## Method 1: Refer to a new task definition

This method of updating a Layer0 application is the easiest to implement, because you do not need to rescale the Service or modify the Loadbalancer. This method is completely transparent to all other components of the application, and using this method does not involve any downtime.

The disadvantage of using this method is that you cannot perform A/B testing of the old and new services, and you cannot control which traffic goes to the old service and which goes to the new one.

**To replace a Deploy to refer to a new task definition:**

1. At the command line, type the following to create a new Deploy: <br />```l0 deploy create [pathToTaskDefinition] [deployName]```<br />Note that if ```[deployName]``` already exists, this step will create a new version of that Deploy.
2. Type the following to update the existing Service: <br />```l0 service update [existingServiceName] [deployName]```<br />By default, the Service you specify in this command will refer to the latest version of ```[deployName]```, if multiple versions of the Deploy exist.<div class="admonition note"><p class="admonition-title">Note</p><br /><p>If you want to refer to a specific version of the Deploy, type the following command instead of the one shown above: <code style="color:white;">l0 service update [serviceName] [deployName]:[deployVersion]</code></p></div>

## Method 2: Create a new Deploy and Service using the same Loadbalancer

This method of updating a Layer0 application is also rather easy to implement. Like the method described in the previous section, this method is completely transparent to all other services and components of the application. This method also you allows you to re-scale the service if necessary, using the ```l0 service scale``` command. Finally, this method allows for indirect A/B testing of the application; you can change the scale of the application, and observe the success and failure rates.

The disadvantage of using this method is that you cannot control the routing of traffic between the old and new versions of the application.

**To create a new Deploy and Service:**

1. At the command line, type the following to create a new Deploy (or a new version of the Deploy, if ```[deployName]``` already exists):<br /> ```l0 deploy create [pathToTaskDefinition] [deployName]```
2. Type the following command to create a new Service that refers to ```[deployName]``` behind an existing Loadbalancer named ```[loadbalancerName]```:<br /> ```l0 service create --loadbalancer [loadbalancerName] [environmentName] [deployName]```
3. Check to make sure that the new Service is working as expected. If it is, and you do not want to keep the old Service, type the following command to delete the old Service: ```l0 service delete [oldServiceName]```

## Method 3: Create a new Deploy, Loadbalancer and Service

The final method of updating a Layer0 service is to create an entirely new Deploy, Loadbalancer and Service. This method gives you complete control over both the new and the old Service, and allows you to perform true A/B testing by routing traffic to individual Services.

The disadvantage of using this method is that you need to implement a method of routing traffic between the new and the old Loadbalancer.

**To create a new Deploy, Loadbalancer and Service:**

1. At the command line, type the following command to create a new Deploy:<br />```l0 deploy create [pathToTaskDefinition] [deployName]```
2. Type the following command to create a new Loadbalancer:<br /> ```l0 loadbalancer create --port [portNumber] [environmentName] [loadbalancerName] [deployName]```<div class="admonition note"><p class="admonition-title">Note</p><br /><p>The value of <code style="color:white;">[loadbalancerName]</code> in the above command must be unique.</p></div>
3. Type the following command to create a new Service: <br />```l0 service create --loadbalancer [loadBalancerName] [environmentName] [serviceName] [deployName]```<div class="admonition note"><p class="admonition-title">Note</p><br /><p>The value of <code style="color:white;">[serviceName]</code> in the above command  must be unique.</p></div>
4. Implement a method of routing traffic between the old and new Services, such as [HAProxy](http://www.haproxy.org) or [Consul](https://www.consul.io).
