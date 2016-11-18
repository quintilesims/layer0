# Shinken
[Shinken](http://www.shinken-monitoring.org/) is an open source monitoring framework for software applications.
XFRA provides [Layer1-Shinken](https://gitlab.imshealth.com/xfra/layer1-shinken/tree/master),
a containerized application that is optimized to integrate into your [Consul-enabled](/guides/consul) Layer0 Environment.

Layer1-Shinken dynamically detects:

* External URLs to query for status
* Layer0 Services registered within Consul that have health checks configured
* Alert notifications are delivered to your team via Slack

## Consul Health Checks
---
A [Consul health check](https://www.consul.io/docs/agent/checks.html) is configurable on each published port in your service by setting environment variables. 
Consul uses these health checks internally to determine the healthy members of a service cluster.

The environment variables are detected by [Registrator](http://gliderlabs.com/registrator/latest/user/backends/#consul), which then automatically populates Consul when your container launches. 
You can specify at most one check per port.


## Consul HTTP Health Check
---
HTTP Checks are the preferred method of monitoring service health in Consul.

The status of the service depends on the HTTP response code: 

| Status | Health |
| - | - |
| 200-299 | Healthy (Green) |
| 429 | Warning (Yellow) | 
| Other / No Response| Unhealthy (Red) |    

HTTP Checks are configured by setting environment variables in your Layer0 Service's task definition. 
Services can configure one health check per port that the container exposes. 
The following environment variables are required for proper health check configuration:

*The "`*`" portion of the environment variable must be replaced with an exposed port on the Service container.*

* **SERVICE_*_NAME**: The name of your Service
* **SERVICE_*_CHECK_HTTP**: The route of your health check
* **SERVICE_*_CHECK_INTERVAL**: The interval between health checks 
* **SERVICE_*_CHECK_TIMEOUT**: The timeout for the health check


Below is a snippet from a task defintion that configures a health check on port 80, against the `/` route, checking every 15 seconds with a 1 second timeout. 
```
"environment": [
    {
        "name": "SERVICE_80_NAME",
        "value": "sample"
    },
    {
        "name": "SERVICE_80_CHECK_HTTP",
        "value": "/"
    },
    {
        "name": "SERVICE_80_CHECK_INTERVAL",
        "value": "15s"
    },
    {
        "name": "SERVICE_80_CHECK_TIMEOUT",
        "value": "1s"
    }
]
```

## Consul TTL Health Check
---
Oftentimes, you'll find your Service does not support HTTP requests, in which case implementing an HTTP endpoint just for health checks is a burden, though a viable option.

An alternative is using TTL (time to live) health checks.
These checks retain their last known state for a given TTL.
The state of the check must be updated periodically over Consul's HTTP interface.
If an external system fails to update the status within a given TTL, the check is set to the failed state.
This mechanism, conceptually similar to a dead man's switch, relies on the application to directly report its health.
For example, a healthy app can periodically PUT a status update to the HTTP endpoint; if the app fails, the TTL will expire and the health check enters a critical state.

To enable a TTL check on Service, simply include the `SERVICE_CHECK_TTL` environment variable with your timeout duration.
For example:
```
"environment": [
    {
        "name": "SERVICE_6379_NAME",
        "value": "db"
    },
    {
        "name": "SERVICE_6379_CHECK_TTL",
        "value": "30s"
    }
]
```

This may sound burndensome as well, but can be quite lightweight. 
XFRA provides a simple example of doing this to provide a health check on Redis [here](https://gitlab.imshealth.com/xfra/redis-ttl/tree/master).


## External Health Checks
---
Layer1-Shinken allows users to configure any number of external HTTP checks.
You can use this to monitor systems not managed by Layer0 from your Layer1-Shinken instance.

External checks simply report on the state of an HTTP(s) URL.
Check health does not affect the availability of your Service, as in the case of [Consul Health Checks](#consul-health-checks), which automatically removes services with failing checks.

To determine the health of a URL, external checks look at the HTTP Status Code returned by the URL:

| Status | Health |
| - | - |
| 200-299 | Healthy (Green) |
| 429 | Warning (Yellow) | 
| Other / No Response| Unhealthy (Red) |

When using external health checks, it is recommended you build a variety of health URLs into your application to monitor the various mission critical aspects that can fail.
This way, you and your team have a simple means of being notified if/when one of those health checks reports a failure.

### Configure External Checks
XFRA provides a Layer0 Service to manage and configure your external checks.
Please review its [README.md](https://gitlab.imshealth.com/xfra/layer1-shinken-ext/blob/master/README.md) for information on configuration options and usage.