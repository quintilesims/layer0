SHELL:=/bin/bash
L0_VERSION:=$(shell git describe --tags | sed "s/v//")
RELEASE:=$(shell cat docs-src/docs/releases.md | sed -n 3p | sed 's/0\.10\../'$(L0_VERSION)'/g')

release:
	$(MAKE) -C api release
	$(MAKE) -C cli release
	$(MAKE) -C runner release
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

github-release:
	$(MAKE) -C cli release
	$(MAKE) -C setup release
	$(MAKE) -C plugins/terraform release

	rm -rf binaries
	cp -R cli/build/ binaries 
	cp -R setup/build/ binaries
	cp -R plugins/terraform/build/  binaries 

	zip -r binaries.zip binaries

update-release:
	# Update Version to Latest and clean up
	sed -i '' 's/0\.10\../'$(L0_VERSION)'/g' README.md docs-src/docs/index.md

	# Add new version to release
	$(shell ex -sc '3i|$(RELEASE)' -cx docs-src/docs/releases.md)

unittest:
	$(MAKE) -C api test
	$(MAKE) -C cli test
	$(MAKE) -C common test
	$(MAKE) -C runner test
	$(MAKE) -C setup test
	$(MAKE) -C plugins/terraform test

smoketest:
	$(MAKE) -C tests/smoke test

systemtest:
	$(MAKE) -C tests/system test

stresstest:
	$(MAKE) -C tests/stress test

install-smoketest:
	$(MAKE) -C cli install-smoketest
	$(MAKE) -C setup install-smoketest
	$(MAKE) -C api deps
	$(MAKE) -C api release
	$(MAKE) -C runner release

apply-smoketest:
	$(MAKE) -C setup apply-smoketest

destroy-smoketest:
	$(MAKE) -C setup destroy-smoketest

.PHONY: release unittest smoketest install-smoketest apply-smoketest destroy-smoketest systemtest benchmark
