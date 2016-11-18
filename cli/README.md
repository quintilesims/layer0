# Layer0 Command Line

The `l0` command line interface (CLI) is used to manage [Layer0](http://docs.xfra.ims.io/) from the shell. 

To get started with Layer0, see [docs.xfra.ims.io](http://docs.xfra.ims.io/). For the `l0` reference guide, see [our CLI guide](http://docs.xfra.ims.io/reference/cli/).

## Layer0 Developers

### Building This Repo

* Run `make build` on Windows, OSX, or Linux from this project's root directory.

### Updating Gomocks

* Run `make all` from `/scripts/update_mocks`. We recommend using the latest version of `gomock` and `mockgen`, which can be installed via:

```
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
```
