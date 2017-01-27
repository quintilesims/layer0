# Layer0 System Tests

### Test Architecture
Each test case contains two components: a subdirectory under `cases/`, and a `<case name>_test.go` file.
The test case directory contains terraform files that outline the resources needed for the test. 
The testing framework will create these resources with `terraform apply` before executing the test, and delete them with `terraform destroy` once the test has completed.

Note that the `cases/modules` directory is not a test case directory.
It holds our [terraform modules](https://www.terraform.io/docs/modules/usage.html) to reduce duplicate code in our terraform files.

If your test case was named `my_simple_example`, you should create the following files and directories:
```
layer0/tests/system
| -- mySimpleExample_test.go
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

func TestMySimpleExample(t *testing.T) {
    c := startSystemTest(t, "cases/my_simple_example", nil)
    defer c.Destroy()
    
    // the rest of your test
}
```

The `startSystemTest` function tells the test to run in parallel, creates a new `framework.SystemTestContext` object using the specified directory (`"cases/my_simple_example"` in this example), and runs `terraform apply`.

In addition to wrapping terraform commands, the `SystemTestContext` object contains some wrappers around the `cli/client` package, e.g.:
```
env := c.GetEnvironment("environment_name")
svc := c.GetService("environment_name", "service_name")
lb := c.GetLoadBalancer("environment_name", "load_balancer_name")
```


### Running a Test
Typically, when writing a new system test you only want to test the one you are working on.
This can be done using the `-run` flag with the name of your test function:
```
go test -run TestMySimpleExample
```

Test verbosity can be increased using the `-v` flag, but this only shows output once the test has completed. 
If you need real time output, use the `-debug` flag and print statements with `logrus.Debugf()`:
```
go test -debug -run TestMySimpleExample
```

When developing a system test, waiting for terraform to setup/destroy the test resources can very time consuming. 
Using the `-dry` flag will tell terraform to use `plan` instead of `apply`. 
Using this method, you can manually create the resources that terraform would make, and then run your test multiple times without terraform destroying and rebuilding the resources:
```
go test -dry -run TestMySimpleExample
```

Some other helpful test flags:
* `-parallel n` - Specifies the number of tests to run in parallel at once
* `-timeout t` - Specifies the timeout for the tests. The default is `10m`, which typically isn't long enough to complete all of the system tests. 


