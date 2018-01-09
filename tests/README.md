# Layer0 Tests
This package contains all non-unit tests for Layer0, including system, smoke, and stress tests. 
The following instructions pertain to system and stress tests. 

### Test Architecture
Each test case contains two components: a subdirectory under `cases/`, and a `<case name>_test.go` file.
The test case directory contains terraform files that outline the resources needed for the test. 
The resources for the test can be created with `terraform apply` before executing the test, and deleted with `terraform destroy` once the test has completed.

Note that the `cases/modules` directory is not a test case directory.
It holds our [terraform modules](https://www.terraform.io/docs/modules/usage.html) to reduce duplicate code in our terraform files.

If your test case was named `my_simple_example`, you should create the following files and directories:
```
layer0/tests/system
|-- mySimpleExample_test.go
|-- cases
|   |-- my_simple_example
|   |    |-- main.tf
```

### Writing a Test
Writing a test is relatively straightforward.
Following the example above, the contents of `mySimpleExample_test.go` should look like:
```go
package system

import (
        "testing"
)

// Test Resources:
// This test creates an environment named 'mse' that has a
// SystemTestService named 'sts'
func TestMySimpleExample(t *testing.T) {
    t.Parallel()

    s := NewSystemTest(t, "cases/my_simple_example", nil)
    s.Terraform.Apply()
    defer s.Terraform.Destroy()
 
    // the rest of your test
}
```

Each test case should have its own file with a matching name. 
The test function has some comments about what resources are created by terraform 
(note: environment names are typically acronyms of the test case, e.g. "my_simple_example" = "mse").

### The SystemTestService
Many system tests use the System Test Service, or [STS](https://github.com/quintilesims/sts).
This is a simple web service whose behavior can be changed through an API.

Checkout the `TestDeadServiceRecreated` test to see how to use the STS client in a test. 

### Running the Tests
To run the system and stress tests, it is required you have the [Layer0 Terraform Plugin](http://layer0.ims.io/reference/terraform-plugin/) in your system's `PATH` variable. 
From the `test/system` or `test/stress` directory, run `make test` to initialize and execute all of the tests.
For more explicit behavior, see the **Test Flags** section below. 

### Test Flags
In addition to the standard `go test` flags, the following have been implemented for system tests:

**-debug** - Test verbosity can be increased using the builtin `-v` flag, but this only shows output once the test has completed. 
If you need real time output, use the `-debug` flag and print statements with `logrus.Debugf()`.

**-dry** - Using the `-dry` flag will swap terraform `apply` and `destroy` commands with `plan`.
When developing a system test, waiting for terraform to setup/destroy the test resources can very time consuming. 
Using this method, you can run your test multiple times without terraform destroying and rebuilding the resources.

Some useful builtin flags:
* `-run nameOfTest` - Executes tests that match the specified name (can be used to run a single test case, or given a param that matches nothing to only run other types of tests).
* `-bench PathOrFileOrTest` - Executes benchmark tests that match the specified identifier.
* `-parallel n` - Specifies the number of tests to run in parallel at once.
* `-short` - Execute the tests in short mode. Long running tests will be skipped.
* `-timeout t` - Specifies the timeout for the tests. 
The default is `10m`, which typically isn't long enough to complete all of the system tests. 

### WARNING: Stress Tests and Service Limits
The stress tests intentionally create and destroy a lot of AWS resources, so you should be aware of the AWS service limits on the account in which these tests will be run.
Stress tests will create configurations of different combinations of Environments, Deploys, Services and Load Balancers.

Currently, the known limits per resource are as follows:

* Environments

Cannot exceed 100 clusters per AWS region.

* Deploys 

If you try to create more than 100 Deploys at a time, you will encounter EOF error messages

* Services

The number of Services you can create will be dependent on the number of Environments and Deploys. The higher the number of Environments and Deploys, the fewer the number of Services.
Generally, you want to create no more than 25 Services at a time.

* Load Balancers

The default limit for AWS Elastic Load Balancing is set to 20 per AWS region.

It is recommended that you use [flow](https://github.com/quintilesims/layer0/blob/develop/scripts/flow.sh) to clean out your layer0 instance prior to running the stress tests.

If you enounter an error running the Stress tests, and you have exceeded the number of allowed Environments, you will see the following error message:

```
AWS Error: clusters cannot have more than 100 elements (code 'InvalidParameterException')
```

If an error like this occurs, you will have to manually enter the AWS console and terminate resources by hand until you're back under the limits.
You may then programatically use Layer0 to destroy the remainder (using [flow](https://github.com/quintilesims/layer0/blob/develop/scripts/flow.sh))

To avoid such a scenario, before you run the stress tests, comment out specific functions that will exceed your service limits, or modify the parameters of some of the stress test functions.

Some references:
* [Viewing account-specific EC2-related service limits](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-resource-limits.html)
* [General ECS service limits](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/service_limits.html)
* [General list of AWS service limits](https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html)
