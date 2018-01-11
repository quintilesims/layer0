init:
	for case in cases/*; do \
		if [ -d "$$case" ] && [ "$$case" != "cases/modules" ]; then \
			pushd "$$case" ; \
			ln -s -f $$(which terraform-provider-layer0) . ; \
                	terraform init ; \
			popd ; \
		fi \
        done

test: init
	go test -v -parallel 10 -timeout 1h

test-short: init
	go test -v -short -parallel 5 -timeout 15m

.PHONY: init test test-short
