# Developer Notes
The following components make up the Layer0 ecosystem.
Each section will describe the component and its responsibilities. 

## Layer0 Setup
The Layer0 Setup tool (commonly called `l0-setup`) is used to provision and manage resources in AWS that a Layer0 instance needs to run correctly.
It is used to "bootstrap" and manage Layer0 instances.

#### Common Functions
*The following commands should be run from the `layer0/setup` directory*

```
# run program
go run main.go

# test program
make test

# build binary
make build

# cross-compile binaries for release
make release
```

## Layer0 CLI
The Layer0 CLI (commonly called `l0`) is used to manage Layer0 resources from the shell.
It interacts with the [Layer0 API](#layer0-api) when running commands. 

#### Common Functions
*The following commands should be run from the `layer0/cli` directory*

```
# run program
go run main.go

# test program
make test

# build binary
make build

# cross-compile binaries for release
make release
```

## Layer0 API
The Layer0 API is a web service that provisions and manages Layer0 resources in AWS.

#### Common Functions
*The following commands should be run from the `layer0/api` directory*

```
# run program
go run main.go

# test program
make test

# build binary
make build

# build and push docker image for release
make release
```

## Layer0 Terraform plugin
The Layer0 Terraform plugin allows Layer0 integration with [Terraform](https://www.terraform.io/). 
Terraform provides documentation on how to create custom providers [here](https://www.terraform.io/guides/writing-custom-terraform-providers.html). 
The plugin is tested through both unit tests and system tests. 

#### Common Functions
*The following commands should be run from the `layer0/plugins/terraform` directory*

```
# run unit tests
make test

# build binary
make build

# build binaries for release
make release
```

## Docs
The documentation for Layer0 exists in `layer0/docs-src`.
We use [mkdocs](http://www.mkdocs.org/) (version >= v0.16.0) to compile our markdown docs into html and css.

#### Common Functions
*The following commands should be run from the `layer0/docs-src` directory*

```
# install dependencies
make deps

# compile docs
make build

# serve docs locally
mkdocs serve
```

## Testing
Any package or subpackage in Layer0 can be tested using the standard `go test` tool.
The smoke tests and system tests require that you have a running Layer0 API to test against. 

#### Common Functions
*The following commands should be run from the `layer0` directory*
```
# run unit tests
make unittest

# run smoke tests
make smoketest

# run system tests
make systemtest
```


# Tools
The following tools are used to help make development with Layer0 easier.

### Gomock
The [Gomock](https://github.com/golang/mock) tool is a mocking framework for Golang.
It uses code generation to create mockable objects.
If any changes are made to mocked interfaces, the code generation will need to run again.

#### Common Functions
*The following commands should be run from the `layer0/scripts` directory*

```
# install dependencies
make deps

# recreate all mocks
make -f Makefile.mocks all

# recreate a subset of mocks
make -f Makefile.mocks <subset>
```

### Go-decorator
The [Go-Decorator](https://github.com/imshealth/go-decorator) tool implements the decorator pattern for Golang.
It uses code generation to create decorated objects.
If any changes are made to decorated interfaces, the code generation will need to run again.

#### Common Functions
*The following commands should be run from the `layer0/scripts` directory*

```
# install dependencies
make deps

# recreate all decorators
make -f Makefile.decorators all
```

### Flow
Flow is a bash script that automates common workflows when developing with Layer0.
This tool requires that the environment variable `LAYER0_PREFIX` is set to the name of you Layer0 instance,
or that the `-p <prefix>` flag is used with any command.

#### Common Functions
*The following commands should be run from the `layer0/scripts` directory*

```
# build the Layer0 api Docker image, push it to Dockerhub, and run the new image as your Layer0 api
flow.sh api

# build the Layer0 runner Docker image, push it to Dockerhub, and run the new image as your Layer0 runner
flow.sh runner

# Delete all entities in your Layer0 (expect for the Layer0 API)
flow.sh delete

# Run all jobs that are in the `PENDING` state (typically only used for local development)
flow.sh runjobs
```

# Development Workflow
A common workflow for working with Layer0

####  Install a Layer0 Instance
This is only required if you don't already have a Layer0 instance.
We will use the Layer0 instances to test our changes AWS.

```
# install a new Layer0
$ cd setup && go run main.go apply <instance>
```

#### Run the API Locally
Now that you have a Layer0 API running in AWS, we can test changes we make against it.
However, it is often much faster to run the Layer0 API locally and test changes against that.
Running the Layer0 API locally requires quite a few environment variables - luckily, those can all be grabbed using the `l0-setup endpoint` command.

```
# get the required environment variables to run your Layer0 API locally using the '-d' flag
$ cd setup && go run main.go endpoint -d <instance>

# set the environment variables returned by the previous command

# run the layer0 api locally
$ go run api/main.go
```

Now, in another shell, you can run your Layer0 CLI against your local Layer0 API by unsetting the `endpoint` variable
```
# unset the endpoint variables
$ unset LAYER0_API_ENDPOINT

# by default, the layer0 cli will use a local api server
# the following command will run against your local api
$ go run cli/main.go environment list
```

#### Test your Changes
Once you have made changes to the code, you should run all of the unit tests:
```
make unittest
```

Once the unit tests are passing, you should run the smoke tests and system tests.
However, these tests require that your updated changes are running in AWS.
This requires building a new docker image for the API and Runner, pushing them to Dockerhub, and updating your Layer0 instance to run the new images.
This can be done using the `flow.sh` script:
```
# use flow to update the API and Runner
./scripts/flow.sh -p <instance> api runner

# note: instead of using the '-p <instance>', you can set a LAYER0_PREFIX environment variable
```

Once your Layer0 has finished updating, you can run the smoketests and system tests
```
# set the endpoint variables
$ l0-setup endpoint <instance>
$ make smoketest
$ make systemtest
```

### Release Process
* Please see [RELEASE.md](RELEASE.md) for instructions on releasing a new version of Layer0.
