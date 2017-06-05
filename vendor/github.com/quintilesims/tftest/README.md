# TFTest

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/quintilesims/tftest/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/quintilesims/tftest)](https://goreportcard.com/report/github.com/quintilesims/tftest)
[![Go Doc](https://godoc.org/github.com/quintilesims/tftest?status.svg)](https://godoc.org/github.com/quintilesims/tftest)

## Overview
TFTest is a tool used for testing application behavior. 
It provides a light wrapper around [Terraform](https://www.terraform.io) to provision and destroy resources for each test case.
Tests are ran as standard [Go tests](https://golang.org/pkg/testing). 

## Getting Started
Check out the [Examples](https://github.com/quintilesims/tftest/tree/master/examples) for some working tests.

The architecture for each test case is simple:
* Create a `TestContext` object
* Create all of the resources for the test case using `TestContext.Apply()`
* Destroy all of the resources for the test case when it completes by using `defer TestContext.Destroy()`
* Run logic on your tests as you would any other test

The following shows a very basic test using TFTest:
```
package main

import (
	"github.com/quintilesims/tftest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
    context := tftest.NewTestContext(t)
    context.Apply()
    defer context.Destroy()

    message := context.Output("message")
    if message != "Hello World" {
        t.Fatalf("Message was '%s', expected 'Hello World'", message)
    }
}
```

## Input Variables
Most test cases will require input variables for Terraform. 
These can be added using the `Var` or `Vars` functions:
```
// insert a single variable
context := tftest.NewTestContext(t, tftest.Var("name", "John"))

// insert multiple variables using Var()
context := tftest.NewTestContext(t, tftest.Var("name", "John"), tftest.Var("age", "35"))

// insert multiple variables using Vars()
context := tftest.NewTestContext(t, tftest.Vars(map[string]string{
	"name": "John", 
	"age":"35",
}))

// add vars on the fly
context := tftest.NewTestContext(t)
context.Vars["name"] = "John"
```

## Output Variables
Outputs from Terraform can be gathered using the `Output` function:
```
context := tftest.NewTestContext(t)
...
name := context.Output("name")
age := context.Output("age")
```

## Context Configuration
#### Execution Directory
By default, a `TestContext` will execute Terraform commands in the `.` directory.
You can change this by using the `Dir` function:
```
context := tftest.NewTestContext(t, tftest.Dir("terraform"))
```

#### Logger
By default, a `TestContext` will log messages using the `t.Log` function. 
You can change this by using the `Log` function:
```
logger := log.New(os.Stdout, "", 0)
context := tftest.NewTestContext(t, tftest.Log(logger))
```

#### Dry Run
It can be very time consuming to spin up/tear down resources while debugging a test. 
Setting the `DryRun` option will execute `terraform plan` in place of `terraform apply`, and `terraform plan -destroy` in place of `terraform destroy`. 
This will allow you to execute the logic of your tests without waiting for the resources to be created and deleted. 

```
context := tftest.NewTestContext(t, tftest.DryRun(true))
```

## Context vs TestContext
The `TestContext` object is a child of the `Context` object. 
It wraps the helper functions of `Context` by calling `t.Fatal` whenever an error occurs. 
If you need to handle errors yourself, or don't have a `testing.T` object (like in [TestMain](https://github.com/quintilesims/tftest/blob/master/examples/setup_teardown/main_test.go)), then use a `Context` object:

```
context := tftest.NewContext()
output, err := context.Apply()
...

testContext := tftest.NewTestContext(t)
context.Apply()
```

# License
This work is published under the MIT license.

Please see the `LICENSE` file for details.
