init:
	pushd module ; \
	ln -s -f $$(which terraform-provider-layer0) . ; \
	terraform init ; \
	popd ; \

test: init
	go test -v -debug -run NOTHING -bench . -timeout 2h 

.PHONY: init test
