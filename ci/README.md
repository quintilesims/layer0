# Layer0 Gitlab CI - Runner

## Build & Push

The XFRA Gitlab CI runner is hosted on `d.ims.io`. To build and push:

```
docker build -t d.ims.io/xfra/gitlab-ci-runner:v1 .
docker push d.ims.io/xfra/gitlab-ci-runner:v1
```

## Further Examples

The Chunnels team has built an extensive CI pipeline, and you can see the various components here:

* [.gitlab-ci.yml](https://gitlab.imshealth.com/tools/ImsHealth.Automation/blob/master/.gitlab-ci.yml)
* [runner](https://gitlab.imshealth.com/tools/gitlab-ci-runner)
* [agent](https://gitlab.imshealth.com/tools/gitlab-ci-agents)
