# Layer0 Setup
Tools for Layer0 installation and management on AWS.

Refer to the official documentation (only available from Office IPs) @ [docs.xfra.ims.io](http://docs.xfra.ims.io)

If the site is down, you can access the source of those docs @ <https://gitlab.imshealth.com/xfra/layer0/tree/master/docs>

## Layer0 Setup - Developer Notes

### Building Layer0 Setup
To build the `l0-setup` binary, run the following command:
```
./make dev
```

### Using a Development API Image
If you'd like to specify a specific version for use with a pushed docker image, set `DEV_VERSION` in ./Makefile accordingly (i.e. "v6.0-rc-1").

### Terraform Version
Versions of l0-setup prior to 0.6.2 use a forked version of terraform at gitlab.imshealth.com/xfra/terraform.
Versions after use the official Terraform binaries at https://www.terraform.io/downloads.html.
Running `make dev` for your current version will ensure you have the proper binaries in place.
