#EC2 Container Registry
ECR is an Amazon implementation of a docker registry.  It acts as a private registry in your AWS account, which can be accessed from any docker client, and Layer0.  Consider using ECR if you have stability issues with hosted docker registries, and do not wish to share your images publicly on [dockerhub](https://hub.docker.com/).

##Setup
When interacting with ECR, you will first need to create a repository and a login to interact from your development machine.

### Repository
Each repository needs to be created by an AWS api call.

```
  > aws ecr create-repository --repository-name myteam/myproject
```

### Login
To authenticate with the ECR service, Amazon provides the `get-login` command, which generates an authentication token, and returns a docker command to set it up

```
  > aws ecr get-login
  # this command will return the following: (password is typically hundreds of characters)
  docker login -u AWS -p password -e none https://aws_account_id.dkr.ecr.us-east-1.amazonaws.com
```
Execute the provided docker command to store the login credentials

Afterward creating the repository and local login credentials you may interact with images (and tags) under this path from a local docker client.

```
  docker pull ${ecr-url}/myteam/myproject
  docker push ${ecr-url}/myteam/myproject:custom-tag-1
```

## Deploy Example
Here we'll walk through using ECR when deploying to Layer0,  Using a very basic wait container.

### Make docker image

Your docker image can be built locally or pulled from dockerhub.  For this example, we made a service that waits and then exits (useful for triggering regular restarts).

> [Dockerfile](https://gitlab.imshealth.com/xfra/production/blob/master/dockerfiles/Dockerfile.wait)
```
FROM busybox

ENV SLEEP_TIME=60

CMD sleep $SLEEP_TIME
```

Then build the file, with the tag `xfra/wait`
```
 > docker build -f Dockerfile.wait -t xfra/wait .
```

### Upload to ECR

After preparing a login and registry, tag the image with the remote url, and use `docker push`

```
  docker tag xfra/wait 111222333444.dkr.ecr.us-east-1.amazonaws.com/xfra-wait
  docker push 111222333444.dkr.ecr.us-east-1.amazonaws.com/xfra-wait
```
!!! note "Note: your account id in this url will be different."

###  Create a deploy

To run this image in Layer0, we create a dockerrun file, describing the instance and any additional variables

> [timeout.Dockerrun.aws.json](https://gitlab.imshealth.com/xfra/production/blob/master/dockerruns/registry.Dockerrun.aws.json)
```
{
  "containerDefinitions": [
    {
      "name": "timeout",
      "image": "111222333444.dkr.ecr.us-east-1.amazonaws.com/xfra-wait:latest",
      "essential": true,
      "memory": 10,
      "environment": [
        { "name": "SLEEP_TIME", "value": "43200" }
      ]
    }
  ]
}
```

And create that in Layer0
```
  l0 deploy create timeout.dockerrun.aws.json timeout
```

### Deploy
Finally, run that deploy as a service or a task. (the service will restart every 12 hours)

```
  l0 service create demo timeoutsvc timeout:latest
```

## References
* [ECR User Guide](http://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_AWSCLI.html)
* [create-repository](http://docs.aws.amazon.com/cli/latest/reference/ecr/create-repository.html)
* [get-login](http://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html)
