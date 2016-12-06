# Release Notes
---
## v0.8.4 | [Download](http://docs.xfra.ims.io)
12.1.2016

#### Features
* Add `terraform-layer0-plugin` see our [docs](http://docs.xfra.ims.io/reference/terraform-plugin) for more information.
* Add terraform plugin examples in "Guestbook" and "Guestbook with RDS" walkthroughs
* Add explicit prefix matching in CLI using '\*' on targets (e.g. `l0 service get d*`)
* Use a single log group in Cloudwatch. 
Deploys are now rendered on `deploy create` to identify log streams instead of on `service/task create`. 
* Add `/admin/config` endpoint for terraform plugin integration
* Remove IMS references and certifcates from `l0-setup`. 
Switch to using a user-defined `dockercfg` file for private registry authentication. 
See [docs](https://docs.xfra.ims.io/setup/update/#upgrading-to-version-084) for more information. 


#### Fixes
* Each environment override on task create will be passed instead of just one
* Tasks that (previously) silently failed to get created will now be inserted into the task scheduler

---
## v0.8.3 | [Download](http://docs.xfra.ims.io)
11.2.2016

#### Features
* Add `--min-count` flag for environment create
* Add `l0 environment setmincount` command
* Report environment instance size

#### Fixes
* Issue hashed ids for environments to prevent duplicate ids on long names

---
## v0.8.2 | [Download](http://docs.xfra.ims.io)
9.26.2016

#### Features
* Add `--no-logs` flag for services and tasks

#### Fixes
* Set max ID length to 12 in order to fix AWS name length constraints

---

## v0.8.1 | [Download](http://docs.xfra.ims.io)
9.9.2016

#### Features
* Add `--user-data` flag for environment create
* Improved release process by downloading Terraform binaries from S3

---
## v0.8.0 | [Download](http://docs.xfra.ims.io)
8.31.2016

#### Features
* Add a `--all` flag for entity list commands:
     * `l0 deploy list --all` will list all versions of each deploy
     * `l0 task list --all` will list tasks which were deleted by the user
* Internal ID management has been standardized
     * The '1' prefix has been removed
* Updated vendoring system to Go 1.6
* Stacktrace logging added to all Layer0 tooling
* Added a `--wait` flag for service create, update, and scale. 
This will wait until the deployment completes before returning.

#### Fixes
* Jobs older than a day are automatically cleaned up
* No longer listing tasks which were deleted by the user
* Internal continuous integration fixes
     * Build scripts and makefiles greatly simplified
* `l0 loadbalancer delete` will wait for IAM role deletion
* Rendered deploys are now tagged properly

------
## v0.7.2 | [Download](http://docs.xfra.ims.io)
7.22.2016

#### Features
* Logs for services, tasks, and jobs have been re-implemented with Cloudwatch
* Environment create takes optional `--size` for EC2 instances in the cluster
* LoadBalancer create can take multiple `--port` flags
* Added new admin debug command to list CLI and API versions
* Jobs older than 24 hours are automatically cleaned up
* Moved the deploy apply command to service update. 
The deploy apply command still exists, but will be depreciated in further releases

#### Fixes
* LoadBalancer delete waits until roles no longer exist before returning
* Increase timeout length during environment delete for autoscaling group
* Commands that run as jobs report their job id instead of giving success message

---
## v0.7.1

6.24.2016 | [Download](http://http://docs.xfra.ims.io/)

#### Features
* API server streams logs to Cloudwatch
* Moved LoadBalancer Delete and Service Delete to the job system
* Added us-west-1 to available regions
* Made "key_pair" a required configuration option for l0-setup
* Updated documentation for cleanly destroying Layer0 instances
* Updated documentation to SSH into Services

#### Fixes

* Jobs no longer get stuck in the "PENDING" state
* Certificate Delete properly removes tags

---
## v0.7.0

5.24.2016 | [Download](http://docs.xfra.ims.io/)

#### Features
* Added a Job system. Long running API tasks now run asynchronously
* Added "--wait" flag for Environment Delete

#### Fixes
* Added proper permissions to delete Certificates
* Allow Services to be created without specifying a LoadBalancer

---
## v0.6.3
5.3.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.3)

#### Features
* Add tags for API Environment, LoadBalancer, and Service

#### Fixes

* Fix RightSizer timeout
* Prevent RightSizer from adjusting API Environment

---
## v0.6.2
4.15.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.2)

*Updating your Layer0 API from previous versions will cause your endpoint URL to change. 
Please review the [update instructions](http://docs.xfra.ims.io/setup/update/#migration-from-v060-or-v061-to-v062) for more information.*

#### Features

* RightSizer added to API 
* Host the API server on ECS

#### Fixes

* Fix long log streams getting cut off
* Nil dereferences in some logs/ps
* Fix service name uniqueness conflict across environments
* DescribeTasks can only handle 100 tasks at once

---
## v0.6.1
4.4.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.1)

#### Features

* 1-Off Tasks are now available in the API and CLI.  For more information on usage consult the documentation on [docs.xfra.ims.io](http://docs.xfra.ims.io/reference/cli/#tasks)

---
## v0.6.0
3.24.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.0)

*This release is incompatible with previous versions of Layer0.
Please see our [update documentation](http://docs.xfra.ims.io/setup/update/) for more information.*

#### Features
* All customer-facing documentation has been updated and is available at [docs.xfra.ims.io](http://docs.xfra.ims.io)
* Services, Environments, and Deploys have been migrated to run in [ECS](https://aws.amazon.com/ecs/)
* CLI command format and output has been significantly updated. Please see the [CLI Reference](http://docs.xfra.ims.io/reference/cli/) for more information
* References to [docker.ims.io](https://docker.ims.io) have been migrated to [d.ims.io](https://d.ims.io)
* Retry logic has been implemented for AWS API calls that get throttled
* The Consul Addon has been removed and is now managed as a normal Layer0 Service.
Please see the [Consul Reference](http://docs.xfra.ims.io/reference/consul/) for more information

---
## v0.6.0-rc.2
3.8.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.0-rc.2)

#### Features
* This is a release candidate for v0.6.0

---
## v0.6.0-rc.1
3.1.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.6.0-rc.1)

#### Features
* This is a release candidate for v0.6.0

---
## v0.5.5
2.24.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.5)

#### Bug Fixes
* Fix consul-agent bootstrapping
* Resolve permissions issue with beanstalk deployment files and S3

---
## v0.5.4
1.26.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.4)

#### Features
* Add documentation for Loadbalancers

#### Bug Fixes
* Fix CLI certificate lookup bug with `loadbalancer:addport`
* CLI no longer prints the Beanstalk loadbalancer url

---
## v0.5.3
1.25.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.3)

#### Features
* Setup locks down API docker image version instead of always using `"l0-api:master"`
* SSL and HTTPS connection terminate at the loadbalancer: SSL -> TCP and HTTPS -> HTTPS

#### Bug Fixes
* Fixed issue where public load balancers would have a mix of public and private subnets
* Bad paths return 404 instead of being redirected to a 200
* Removed defunct `certificate:apply` command

---
## v0.5.2
1.15.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.2)

#### Bug Fixes
* CLI `loadbalancer:update` command now has correct url path

---
## v0.5.1
1.14.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.1)

#### Bug Fixes
* Setup now sets `iam:PassRole` permission correctly

---
## v0.5.0
1.13.2016 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.5.0)

#### Feautures
* Loadbalancers as a first-class entity. Please see [Usage](https://gitlab.imshealth.com/xfra/l0-cli/tree/master#create-a-load-balancer) for more information
* Setup now automatically initializes tag database during `apply`;
users no longer have to manually use the `admin/sql` endpoint in swagger.
* The `admin/sql` command had been added to the CLI as `admin:sql`

#### Bug Fixes
* CLI `list` commands return multiple entities if matches exist
* Layer0 Consul agent always restarts

---
# Layer0 v0.5 is now available!
---
## v0.4.2
12.8.2015 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.4.2)

#### Feautures
* Layer0 can be installed into an existing VPC.
Please see [Usage](https://github.com/quintilesims/layer0/tree/master/setup#install-into-an-existing-vpc) for more information
* Random authentication token generation for API
* Expanded CLI usage/help information

#### Bug Fixes
* CLI returns non-zero exit-codes on errors
* CLI (and API) will now consider an exact name match before prefixes
    * Your environments `test` and `test2` can now live in harmony
* Setup gives errors when Layer0 prefix is > 15 characters (rather than truncating)

---
## v0.4.1
11.18.2015 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.4.1)

#### Feautures
* Custom authentication option for API. Please see [Usage](https://github.com/quintilesims/layer0/tree/master/setup#custom-authentication) for more information
* Updated Setup now uses `Terraform` for all infrastructure management
* Improved error messages in CLI
* CLI only shows the latest version of each deploy with `deploy` and `deploy:get` commands

#### Bug Fixes
* Conflicting AWS and Layer0 terminology in error messages
* Migrate docker.appature.com references to docker.ims.io
* Missing rds and autoscaling permissions in bootstrap user policy

---
## v0.4.0
1.14.2015 | [Download](https://github.com/quintilesims/layer0-release/tree/master/v0.4.0)

*This release has changed some of the versions previously released to be more compliant with http://semver.org.  For more information please read [Reversioning](https://github.com/quintilesims/layer0-release/blob/master/REVERSIONING.md).*

#### Feautures
* Tagging for Layer0 entities. Please see [Usage](https://gitlab.imshealth.com/xfra/l0-cli/tree/master#aside-fuzzy-match) for more information
* Direct access to logs in the CLI using `service:logs`
* Access to API logs in the CLI using `admin:logs`
* Default authentication for deploys. Please see [Usage](https://gitlab.imshealth.com/xfra/l0-cli/tree/master#default-auth) for more information

#### Bug Fixes
* Issues with the Consul provisioning process
* Deletion not working for Services
* Deletion not working for Consul
