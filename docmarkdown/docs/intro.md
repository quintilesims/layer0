# Layer0 Introduction

In recent years, the process of deploying applications has seen incredible innovation. However, this innovation has taken a somewhat simple task and made it into something quite [complicated](https://www.nginx.com/blog/microservices-at-netflix-architectural-best-practices/). Cloud providers, load balancing, virtual servers, IP subnets, and a continuing list of technological considerations are not only required to be understood, but their creation and management must be automated for a modern application to be successful at scale.

The burden of understanding a complicated and ever-growing infrastructure is a large aspect of what Layer0 is trying to fix. We've already done the leg work for huge swathes of your backend infrastructure, and we've made it easy to tear down and start over again, too. Meanwhile, you can develop locally using [Docker](https://docs.docker.com/engine/understanding-docker/) and be assured that your application will properly translate to the cloud when you're ready to deploy.

Layer0 requires a solid understanding of Docker to get the most out of it. We highly recommend starting with [Docker's Understanding the Architecture](https://docs.docker.com/engine/understanding-docker/) to learn more about using Docker locally and in the cloud. We also recommend the [Twelve-Factor App](http://12factor.net/) primer, which is a critical resource for understanding how to build a microservice.

---
## Layer0 Concepts

The following concepts are core Layer0 abstractions for the technologies and features we use [behind the scenes](reference/architecture.md). These terms will be used throughout our guides, so having a general understanding of them is helpful.

### Certificates

An SSL certificate obtained from a valid [Certificate Authority (CA)](http://wikit.rxcorp.com/index.php/Production_Dns_Request#CSR). You can use these certificates to secure your HTTPS services by applying them to your Layer0 load balancers.

### Deploys

A [multicontainer docker configuration](http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/create_deploy_docker_v2config.html). This configuration file details how to deploy your application. We have several [sample applications](https://gitlab.imshealth.com/xfra/layer0-samples) available that show what these files look like --- they're called `Dockerrun.aws.json` within each sample app.

### Load Balancers

A powerful tool that gives you the basic building blocks for high-availability, scaling, and HTTPS. We currently use Amazon's [Elastic Load Balancing](https://aws.amazon.com/elasticloadbalancing/), and it pays to understand the basics of this service when working with Layer0.

### Services

Your running Layer0 application. We also use the term `service` for tools such as Consul, Logstash, and Shinken because they are Layer0 applications that we've pre-built for you.

### Environments

A logical grouping of services. Typically, you would make a single environment for each tier of your application, such as `dev`, `staging`, and `prod`.
