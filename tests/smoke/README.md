# Layer0 Smoke Tests

These tests use the [Bats](https://github.com/sstephenson/bats) framework to run. 
Follow the installation instructions in the provided link. 

#### Local Config
Environment Variables must be populated with the contents of `l0-setup endpoint -d <prefix>`

#### Running

Make sure you have cleared out any existing environments, services, etc. from your Layer0 instance before running the test. 
You can use `flow delete` to accomplish this.

From the `layer0` directory, run `make smoketest`

Or, from the `layer0/tests/smoke` directory, run `make test`

#### Tips and Tricks

* Leave no trace - delete any resources that were created during the test
* Resource deletion typically runs asynchronous. 
Use the `--wait` flag to ensure the test doesn't continue until the resource has been deleted
* Place any non `.bats` files required for your test into the `common` directory
