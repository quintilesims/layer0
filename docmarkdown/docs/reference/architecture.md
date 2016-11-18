# Layer0 Architecture

Layer0 is built on top of the following primary technologies:

* Application Container: [Docker](https://docs.docker.com/engine/understanding-docker/)
* Cloud Provider: [Amazon Web Services](https://aws.amazon.com/)
  * Provider Billing, Configuration: [Bring Your Own Ops (BYOO)](http://wikit.rxcorp.com/index.php/BYOO)
* Container Management: [Amazon EC2 Container Service (ECS)](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html)
* Load Balancing: [Amazon Elastic Load Balancing](https://aws.amazon.com/elasticloadbalancing/)
* Infrastructure Configuration: Hashicorp [Terraform](https://www.terraform.io/docs/index.html)
* Identity Management: [Auth0](https://auth0.com/docs) and IMS Active Directory
