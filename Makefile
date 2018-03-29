SHELL:=/bin/bash
L0_VERSION?=$(shell git describe --tags)

build: 
	$(MAKE) -C api build        
	$(MAKE) -C setup build       
	$(MAKE) -C cli build      

install:
	$(MAKE) -C cli install
	$(MAKE) -C setup install

release:
	$(MAKE) -C api release
	$(MAKE) -C cli release
	$(MAKE) -C setup release
	$(MAKE) -C plugins/terraform release

	rm -rf build
	for os in linux darwin windows; do \
		cp -R cli/build . ; \
		cp -R setup/build . ; \
		cp -R plugins/terraform/build . ; \
		cd build/$$os && zip -r layer0_$(L0_VERSION)_$$os.zip * && cd ../.. ; \
		aws s3 cp build/$$os/layer0_$(L0_VERSION)_$$os.zip s3://xfra-layer0/release/$(L0_VERSION)/layer0_$(L0_VERSION)_$$os.zip ; \
	done

unittest:
	$(MAKE) -C api test
	$(MAKE) -C cli test
	$(MAKE) -C client test
	$(MAKE) -C common test
	$(MAKE) -C setup test
	$(MAKE) -C plugins/terraform test

smoketest:
	$(MAKE) -C tests/smoke test

systemtest:
	 $(MAKE) -C tests/system test

.PHONY: release unittest smoketest install-smoketest apply-smoketest destroy-smoketest systemtest
