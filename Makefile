SHELL:=/bin/bash

release:
	$(MAKE) -C api release
	$(MAKE) -C cli release
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

stresstest:
	$(MAKE) -C tests/stress test

systemtest:
	$(MAKE) -C tests/system test

.PHONY: release unittest smoketest stresstest systemtest
