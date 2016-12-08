# Task Definitions
This guide gives some overview into the composition of a task definition.
For more comprehensive documentation, we recommend taking a look at the official AWS docs:

* [Creating a Task Definition](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/create-task-definition.html)
* [Task Definition Parameters](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html)


## Sample
The following snippet contains the task definition for the [Guestbook](/guides/guestbook) application
```
{
    "AWSEBDockerrunVersion": 2,
    "containerDefinitions": [
        {
            "name": "l0-demo-guestbook",
            "image": "d.ims.io/xfra/l0-guestbook",
            "essential": true,
            "memory": 128,
            "portMappings": [
                {
                    "hostPort": 80,
                    "containerPort": 80
                }
            ],
            "environment": [
                {
                    "name": "SERVICE_80_NAME",
                    "value": "l0-guestbook"
                }
            ]
        }
    ]
}
```

* **Name** The name of the container

!!! warning
If you wish to update your task definition, the container names **must** remain the same.
If any container names are changed or removed in an updated task definition,
ECS will not know how the existing container(s) should be mapped over and you will not be able to deploy the updated task definition.
If you encounter a scenario where you must change or remove a container's name in a task definition, we recommend re-creating the Layer0 Deploy and Service.


* **Image** The Docker image used to build the container. The image format is `url/image:tag`
    * The `url` specifies which Docker Repo to pull the image from (e.g. `d.ims.io`).
If `url` is not specified, [Docker Hub](https://hub.docker.com/) is used
    * The `image` specifies the name of the image to grab
    * The `tag` specifies which version of image to grab.
If `tag` is not specified, `:latest` is used
* **Essential** If set to `true`, all other containers in the task definition will be stopped if that container fails or stops for any reason.
Otherwise, the container's failure will not affect the rest of the containers in the task definition.
* **Memory** The number of MiB of memory to reserve for the container.
If your container attempts to exceed the memory allocated here, the container is killed
* **PortMappings** A list of hostPort, containerPort mappings for the container
    * **HostPort** The port number on the host instance reserved for your container.
If your Layer0 Service is behind a Layer0 Load Balancer, this should map to an `instancePort` on the Layer0 Load Balancer.
    * **ContainerPort** The port number the container should receive traffic on.
Any traffic received from the instance's `hostPort` will be forwarded to the container on this port
* **Environment** A list of key/value pairs that will be available to the container as environment variables
