SHELL:=/bin/bash
L0_VERSION:=$(shell git describe --tags)

release:
	$(MAKE) -C api release
	$(MAKE) -C cli release
	$(MAKE) -C docs-src release
	$(MAKE) -C runner release
	$(MAKE) -C setup release
	$(MAKE) -C plugins/terraform release

	rm -rf build
	for os in Linux macOS Windows; do \
		cp -R cli/build . ; \
		cp -R setup/build . ; \
		cp -R plugins/terraform/build . ; \
		cd build/ && zip -r $$os.zip $$os && cd .. ; \
	done

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
