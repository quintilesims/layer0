# Layer0

Layer0 is a framework that helps you deploy your docker container to the cloud with minimal fuss. Using a simple command line interface (CLI), you can manage the entire life cycle of your application without having to focus on infrastructure.

To get started with Layer0, see [quintilesims.github.io/layer0](https://quintilesims.github.io/layer0/).

# Components
The following components comprise the Layer0 Application. 

## Layer0 Setup
The Layer0 Setup tool (commonly called `l0-setup`) is used to provision and manage resources in AWS that a Layer0 instance needs to run correctly. 

#### Common Functions
*The following commands should be ran from the `layer0/setup` directory*

| Function      | Command       |
| ------------- |:--------------|
| Run program   | `go run main.go` |
| Build binary  | `go build -o l0-setup` |
| Cross-Compile | `make build` |
| Release       | `make release` |

## Layer0 CLI
The Layer0 CLI (commonly called `l0`) is used to manage Layer0 resources from the shell.

#### Common Functions
*The following commands should be ran from the `layer0/cli` directory*

| Function      | Command       |
| ------------- |:--------------|
| Run program   | `go run main.go` |
| Build binary  | `go build -o l0` |
| Cross-Compile | `make build` |
| Release       | `make release` |

## Layer0 API
The Layer0 API is a web service that provisions and manages Layer0 resources in AWS.

#### Common Functions
*The following commands should be ran from the `layer0/api` directory*

| Function      | Command       |
| ------------- |:--------------|
| Run program        | `go run main.go` |
| Build binary       | `go build -o l0-api` |
| Build Docker Image | `make build` |
| Release            | `make release` |

## Layer0 Runner
The Layer0 Runner is a service that runs jobs created by the Layer0 API.

#### Common Functions
*The following commands should be ran from the `layer0/runner` directory*

| Function      | Command       |
| ------------- |:--------------|
| Run program        | `go run main.go` |
| Build binary       | `go build -o l0-runner` |
| Build Docker Image | `make build` |
| Release            | `make release` |

## Docs
The documentation for Layer0 exists in `layer0/docs-src`. 
We use [mkdocs](http://www.mkdocs.org/) to compile our markdown docs into html and css.

#### Common Functions
*The following commands should be ran from the `layer0/docs-src` directory*

| Function      | Command       |
| ------------- |:--------------|
| Install Deps       | `make deps` |
| Compile Docs       | `make build` |
| Serve Docs Locally | `mkdocs serve` |

# Testing
Any package or subpackage in Layer0 can be tested using the standard `go test` tool.
The smoke tests and system tests require that you have a running Layer0 API server with the proper environment variables in place (which can be gathered from the `l0-setup endpoint` command).

#### Common Functions
*The following commands should be ran from the `layer0` directory*

| Function      | Command       |
| ------------- |:--------------|
| Run Unit tests       | `make unittest` |
| Run Smoke Tests      | `make smoketest` |
| Run System Tests     | `make systemtest` |


# Tools 
The following tools are used to help make development with Layer0 easier.

## Gomock
The [Gomock](https://github.com/golang/mock) tool is a mocking framework for Golang. 
It uses code generation to create mockable objects. 
If any changes are made to mocked interfaces, the code generation will need to run again.

#### Common Functions
*The following commands should be ran from the `layer0/scripts` directory*

| Function      | Command       |
| ------------- |:--------------|
| Install Deps               | `make -f Makefile.mocks deps` |
| Recreate all mocks         | `make -f Makefile.mocks all` |
| Recreate a subset of mocks | `make -f Makefile.mocks <subset>` |

## Go-decorator
The [Go-Decorator](https://github.com/imshealth/go-decorator) tool implements the decorator pattern for Golang.
It uses code generation to create decorated objects. 
If any changes are made to decorated interfaces, the code generation will need to run again.

#### Common Functions
*The following commands should be ran from the `layer0/scripts` directory*

| Function      | Command       |
| ------------- |:--------------|
| Install Deps                    | `make -f Makefile.decorators deps` |
| Recreate all decorators         | `make -f Makefile.decorators all` |
| Recreate a subset of decorators | `make -f Makefile.decorators <subset>` |

## Flow
Flow is a bash script that automates common workflows when developing with Layer0. 
This tool requires to have an environment variable `LAYER0_PREFIX` set to the name of you Layer0 instance, 
or that the `-p <prefix>` flag is used with any command.

#### Common Functions
*The following commands should be ran from the `layer0/scripts` directory*

| Function      | Command       |
| ------------- |:--------------|
| Build the API Docker image, push it to Dockerhub, and run the new image as your Layer0 API | `flow.sh api` |
| Build the Runner Docker image, push it to Dockerhub, and run the new image as your Layer0 Runner | `flow.sh runner`|
| Delete all entities in your Layer0 that aren't the API | `flow.sh delete` |
| Run all jobs that are in the `PENDING` state (only for use with local development) | `flow.sh runjobs` |

# Development Workflow
A common workflow for working with Layer0

####  Install a Layer0 Instance
This is only required if you don't already have a Layer0 instance.
We will use the Layer0 instances to test our changes AWS.

```
# install a new Layer0
$ go run setup/main.go apply <instance>
```

#### Run the API Locally
Now that you have a Layer0 API running in AWS, we can test changes we make against it. 
However, it is often much faster to run the Layer0 API locally and test changes against that. 
Running the Layer0 API locally requires quite a few environment variables - luckily, those can all be grabbed using the `l0-setup endpoint` command.

```
# get the required environment variables to run your Layer0 API locally using the '-d' flag
$ go run setup/main.go endpoint -d <instance>

# set the environment variables returned by the previous command
# run the layer0 api locally
$ go run api/main.go
```

Now, in another shell, you can run your Layer0 CLI against your local Layer0 API by unsetting the `endpoint` variables
```
# unset the endpoint variables
$ unset LAYER0_API_ENDPOINT
$ unset LAYER0_AUTH_TOKEN

# by default, the layer0 cli will use a local api server
# the following command will run against your local api
$ go run cli/main.go environment list
```

#### Test your Changes
Once you have make changes to your code, you should run all of the unit tests:
```
make unittest
```

Once the unit tests are passing, you should run the smoke tests and system tests. 
However, these tests require that your update changes are running in AWS. 
This requires bubilding a new docker image for the API and Runner, pushing them to github, and updating your Layer0 instance to run the new images. 
This can be done using the `flow.sh` script
```
./scripts/flow.sh -p <instance> api runner
```

Once your Layer0 has finished running, you can run the smoketests and system tests
```
# set the endpoint variables
$ l0-setup endpoint <instance>
$ make smoketest
$ make systemtest
```

### Release Process
* Please see [RELEASE.md](RELEASE.md) for instructions on releasing a new version of Layer0.

