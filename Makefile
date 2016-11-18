L0_VERSION?=$(shell git describe --tags)

release:
	$(MAKE) -C api release
	$(MAKE) -C cli release
	$(MAKE) -C runner release
	$(MAKE) -C setup release

	rm -rf build
	for os in linux darwin windows; do \
		cp -R cli/build . ; \
		cp -R setup/build . ; \
		cd build/$$os && zip -r layer0_$(L0_VERSION)_$$os.zip * && cd ../.. ; \
		aws s3 cp build/$$os/layer0_$(L0_VERSION)_$$os.zip s3://xfra-layer0/release/$(L0_VERSION)/layer0_$(L0_VERSION)_$$os.zip ; \
	done

unittest:
	$(MAKE) -C api test
	$(MAKE) -C cli test
	$(MAKE) -C runner test
	$(MAKE) -C setup test

.PHONY: release unittest
