# Updating a Layer0 service
There are three methods of updating an existing Layer0 service. The first method is to update the existing Deploy to refer to a new Docker task definition. The second method is to create a new Service that uses the same Loadbalancer. The third method is to create both a new Loadbalancer and a new Service.

There are advantages and disadvantages to each of these methods. The following sections discuss the advantages and disadvantages of using each method, and include procedures for implementing each method.

## Method 1: Refer to a new task definition

This method of updating a Layer0 application is the easiest to implement, because you do not need to rescale the Service or modify the Loadbalancer. This method is completely transparent to all other components of the application, and using this method does not involve any downtime.

The disadvantage of using this method is that you cannot perform A/B testing of the old and new services, and you cannot control which traffic goes to the old service and which goes to the new one.

**To replace a Deploy to refer to a new task definition:**

At the command line, type the following to create a new Deploy:

```
l0 deploy create taskDefPath deployName
```

`taskDefPath` is the path to the ECS Task Definition. Note that if `deployName` already exists, this step will create a new version of that Deploy.

Use [l0 service update](cli/#service-update) to update the existing service:
   
```
l0 service update serviceName deployName[:deployVersion]
```

By default, the service name you specify in this command will refer to the latest version of `deployName`. You can optionally specify a specific version of the deploy, as shown above.

## Method 2: Create a new Deploy and Service using the same Loadbalancer

This method of updating a Layer0 application is also rather easy to implement. Like the method described in the previous section, this method is completely transparent to all other services and components of the application. This method also you allows you to re-scale the service if necessary, using the [l0 service scale](cli/#service-scale) command. Finally, this method allows for indirect A/B testing of the application; you can change the scale of the application, and observe the success and failure rates.

The disadvantage of using this method is that you cannot control the routing of traffic between the old and new versions of the application.

**To create a new Deploy and Service:**

At the command line, type the following to create a new deploy or a new version of a deploy:

```
l0 deploy create taskDefPath deployName
```

`taskDefPath` is the path to the ECS Task Definition. Note that if `deployName` already exists, this step will create a new version of that Deploy.

Use [l0 service create](cli/#service-create) to create a new service that uses `deployName` behind an existing load balancer named `loadBalancerName`

```
l0 service create --loadbalancer [environmentName:]loadBalancerName environmentName serviceName deployName[:deployVersion]
```

By default, the service name you specify in this command will refer to the latest version of `deployName`. You can optionally specify a specific version of the deploy, as shown above. You can also optionally specify the name of the environment, `environmentName` where the load balancer exists. 

Check to make sure that the new service is working as expected. If it is, and you do not want to keep the old service, delete the old service: 

```
l0 service delete service
```

## Method 3: Create a new Deploy, Loadbalancer and Service

The final method of updating a Layer0 service is to create an entirely new Deploy, Load Balancer and Service. This method gives you complete control over both the new and the old Service, and allows you to perform true A/B testing by routing traffic to individual Services.

The disadvantage of using this method is that you need to implement a method of routing traffic between the new and the old Load Balancer.

**To create a new Deploy, Load Balancer and Service:**

Type the following command to create a new Deploy:

```
l0 deploy create taskDefPath deployName
```

`taskDefPath` is the path to the ECS Task Definition. Note that if `deployName` already exists, this step will create a new version of that Deploy.

Use [l0 loadbalancer create](cli/#loadbalancer-create) to create a new Load Balancer:

```
l0 loadbalancer create --port port environmentName loadBalancerName deployName
```

* `port` is the port configuration for the listener of the Load Balancer. Valid pattern is `hostPort:containerPort/protocol`. Multiple ports can be specified using `--port port1 --port port2 ...`.
    * `hostPort` - The port that the load balancer will listen for traffic on.
    * `containerPort` - The port that the load balancer will forward traffic to.
    * `protocol` - The protocol to use when forwarding traffic (acceptable values: TCP, SSL, HTTP, and HTTPS).

!!! note
    The value of `loadbalancerName` in the above command must be unique to the Environment.

Use [l0 service create](cli/#service-create) to create a new Service using the Load Balancer you just created: 

```
l0 service create --loadbalancer loadBalancerName environmentName serviceName deployName
```

!!! note
    The value of `serviceName` in the above command  must be unique to the Environment.

Implement a method of routing traffic between the old and new Services, such as [HAProxy](http://www.haproxy.org) or [Consul](https://www.consul.io).
