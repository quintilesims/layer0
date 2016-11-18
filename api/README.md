# Layer0 API Service

Refer to the official documentation (only available from Office IPs) @ [docs.xfra.ims.io](http://docs.xfra.ims.io)

## Developer Notes

### Local Setup

For developing patches to Layer0, it can be helpful to run the api
locally.  To get started, follow the [Layer0 Setup guide](http://docs.xfra.ims.io/setup/install/#install-layer0), and export
all the [VPC configuration](http://docs.xfra.ims.io/setup/install/#install-layer0) locally (see the Config section below).

#### Build
1. Clone this repository
2. `$ cd api`
3. `$ make build`

#### Config

Export the following to your shell of choice. Values for the variables prefixed with **LAYER0_** are generated from [Layer0 Setup](http://docs.xfra.ims.io/setup/install/#configure-layer0).

```
export LAYER0_AWS_ACCESS_KEY_ID=<access key>
export LAYER0_AWS_SECRET_ACCESS_KEY=<secret key>
export LAYER0_PREFIX="<your custom prefix>"
export LAYER0_PUBLIC_SUBNETS="<public subnets>"
export LAYER0_PRIVATE_SUBNETS="<private subnets>"
export LAYER0_ECS_INSTANCE_PROFILE="<ecs role instanceprofile>"
export LAYER0_VPCID="<your vpc>"
export LAYER0_S3_BUCKET="<bucket name>"
```

Notes:
* `LAYER0_ECS_INSTANCE_PROFILE` must be an [Instance Profile](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_switch-role-ec2_instance-profiles.html) name.
* The IAM user roles are created by the `l0-setup` tool.  Examine `*_policy.json` under `setup\terraform` for the details.

##### Optional
* set `LAYER0_AWS_REGION` if your instance is outside of `us-west-2`
* set `LAYER0_PORT` to a port suffix to override the default of `9090`

### Database
The api relies on a sql backend for storing instance tags.  As a
default, the api server will use sqlite (in memory).  Sqlite requires
no setup, but disappears on restart.  For more permanence, use a local
mysql instance or connect to the RDS cluster prepared by `l0-setup`.

##### Sqlite
* No further configuration
* restarting layer0 forgets existing tags (existing services and environments become inaccessible)

##### Local Mysql
* Create a granting user, and select values for user, password, and database for this layer0 instance.
* Add these settings to your environment before running the `api`.

```
export LAYER0_MYSQL_CONNECTION="<user>:<password>@tcp(localhost:3306)/<database>"
export LAYER0_MYSQL_ADMIN_CONNECTION="<master_user>:<master_password>@tcp(localhost:3306)/"
```

##### Amazon RDS
To use the database created by `l0-setup` from your desktop, you'll need to first create a tunnel to access it.

* You can create a tunnel into your layer0 VPC with the ssh sample
* Specify a local tunnel, for example -L 13306:<rds-endpoint>:3306
* Prepare your layer0 config to talk to the local tunnel.

```
export LAYER0_MYSQL_CONNECTION="<user>:<password>@tcp(localhost:13306)/<database>"
export LAYER0_MYSQL_ADMIN_CONNECTION="<master_user>:<master_password>@tcp(localhost:13306)/"
```

#### Database Init
For either of the sql options, a one time setup is needed which can be triggered from the apidocs UI.

* To prepare the database, start the service and navigate to http://localhost:9090/apidocs/#!/admin/UpdateSql
* enter `{"Version":"Latest"}` as the body and **Try it out!**
* On success, the response will be 204 (no content)
* Also, once the database is prepared, [the get api](http://localhost:9090/apidocs/#!/admin/GetSql) will return debug information like table structure.

### Run Layer0 API
```
$ ./api
```

#### API Output
When the service starts successfully, you'll see log lines like this:
```
$ ./api
2015-10-26T11:52:06-07:00 [INFO]: l0-api v0.4.0-58-g11418cf
[restful] 2015/10/26 11:52:06 log.go:30: [restful/swagger] listing is available at /apidocs.json
[restful] 2015/10/26 11:52:06 log.go:30: [restful/swagger] /apidocs/ is mapped to folder api/external/swagger-ui/dist
2015-10-26T11:52:06-07:00 [INFO]: Service on localhost:9090
```

Alternatively, a failure (typically missing environment settings) will be reported to stdout and the api service will terminate.

### Accessing the API Locally
Visit [localhost:9090/apidocs/](http://localhost:9090/apidocs) to explore the
API. Viewer powered by [swagger-ui](https://github.com/swagger-api/swagger-ui).

Also you can `set LAYER0_API_ENDPOINT` to address your local instance from the `l0` command line.

```
$ export LAYER0_API_ENDPOINT=http://localhost:9090
```
*(or change the port as appropriate)*

## Updating Gomocks

* Run `make all` from `/scripts/update_mocks`. We recommend using the latest version of `gomock` and `mockgen`, which can be installed via:

```
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
```
