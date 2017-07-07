# Release Notes
---
## v0.10.1 | [Download](http://layer0.ims.io/#download)
7.7.2017

#### Fixes
* Includes a hotfix for "Pending Task" issue

---
## v0.10.0 | [Download](http://layer0.ims.io/#download)
6.12.2017

#### Features
* Windows container support
* Linked environments
* No limit for --copies option on task create
* Allow --ami specification on environment create
* New walkthroughs and examples

#### Fixes
* Refactor of l0-setup; use terraform modules for installation
* Improved performance of DynamoDB tag and job tables
* Documentation consistency issues

---
## v0.9.0 | [Download](http://layer0.ims.io/releases/)
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

## v0.8.4 | [Download](http://layer0.ims.io/releases/)
12.1.2016

#### Features
* Add `terraform-layer0-plugin`
* Add terraform plugin examples in "Guestbook" and "Guestbook with RDS" walkthroughs
* Add explicit prefix matching in CLI using '\*' on targets (e.g. `l0 service get d*`)
* Use a single log group in Cloudwatch. 
Deploys are now rendered on `deploy create` to identify log streams instead of on `service/task create`. 
* Add `/admin/config` endpoint for terraform plugin integration
* Remove IMS references and certifcates from `l0-setup`. 
Switch to using a user-defined `dockercfg` file for private registry authentication. 


#### Fixes
* Each environment override on task create will be passed instead of just one
* Tasks that (previously) silently failed to get created will now be inserted into the task scheduler

