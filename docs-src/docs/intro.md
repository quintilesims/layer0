# Layer0 Introduction

In recent years, the process of deploying applications has seen incredible innovation. However, this innovation has taken a somewhat simple task and made it into something quite [complicated](https://www.nginx.com/blog/microservices-at-netflix-architectural-best-practices/). Cloud providers, load balancing, virtual servers, IP subnets, and a continuing list of technological considerations are not only required to be understood, but their creation and management must be automated for a modern application to be successful at scale.

The burden of understanding a complicated and ever-growing infrastructure is a large aspect of what Layer0 is trying to fix. We've already done the leg work for huge swathes of your backend infrastructure, and we've made it easy to tear down and start over again, too. Meanwhile, you can develop locally using [Docker](https://docs.docker.com/engine/understanding-docker/) and be assured that your application will properly translate to the cloud when you're ready to deploy.

Layer0 requires a solid understanding of Docker to get the most out of it. We highly recommend starting with [Docker's Understanding the Architecture](https://docs.docker.com/engine/understanding-docker/) to learn more about using Docker locally and in the cloud. We also recommend the [Twelve-Factor App](http://12factor.net/) primer, which is a critical resource for understanding how to build a microservice.

---
## Layer0 Concepts

The following concepts are core Layer0 abstractions for the technologies and features we use [behind the scenes](reference/architecture.md). These terms will be used throughout our guides, so having a general understanding of them is helpful.

### Certificates

SSL certificates obtained from a valid [Certificate Authority (CA)](https://en.wikipedia.org/wiki/Certificate_authority). You can use these certificates to secure your HTTPS services by applying them to your Layer0 load balancers.

### Deploys

[ECS Task Definitions](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_defintions.html). These configuration files detail how to deploy your application. We have several [sample applications](https://github.com/quintilesims/guides) available that show what these files look like --- they're called `Dockerrun.aws.json` within each sample app.

### Tasks

Manual one-off commands that don't necessarily make sense to keep running, or to restart when they finish. These run using Amazon's `RunTask` action (more info [here](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduling_tasks.html)), and are "ideally suited for processes such as batch jobs that perform work and then stop."

### Load Balancers

Powerful tools that give you the basic building blocks for high-availability, scaling, and HTTPS. We currently use Amazon's [Elastic Load Balancing](https://aws.amazon.com/elasticloadbalancing/), and it pays to understand the basics of this service when working with Layer0.

### Services

Your running Layer0 applications. We also use the term `service` for tools such as Consul, for which we provide a pre-built [sample implementation](guides/consul) using Layer0.

### Environments

Logical groupings of services. Typically, you would make a single environment for each tier of your application, such as `dev`, `staging`, and `prod`. Additionally an environment can be either static or dynamic. Static environments should be used when you require more fine grained control over an EC2 instance that a container will run on. Dynamic environments (the default) allow you to run containers without having to worry about managing and scaling clusters of EC2 instances.
