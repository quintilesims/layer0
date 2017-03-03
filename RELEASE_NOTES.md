# Release Notes
---
## v0.9.0 | [Download](https://quintilesims.github.io/layer0/)
3.3.2017

#### Features
* Task and service [placement contraints](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-placement-constraints.html).
* Elastic Load Balancer (ELB) healthcheck configuration.
* System test framework.
* End-to-end smoketest automation.

#### Fixes
* Database layer refactor.
* l0-setup allows for .dockercfg or config.json auth config(s).
* Bulk entity list improvements.
* Many, many documentation updates.

#### Deprecated
* Layer0 certificates.

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

