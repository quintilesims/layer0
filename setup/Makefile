L0_VERSION?=$(shell git describe --tags)
TERRAFORM_VERSION=0.8.5

build-directories:
	rm -rf build
	for os in linux darwin windows; do \
		mkdir -p build/$$os/bin ; \
		cp *.template build/$$os/ ; \
		cp -r templates build/$$os/templates ; \
	done

download-terraform: build-directories
	for os in linux darwin windows; do \
		cd build/$$os/bin ; \
		aws s3 cp s3://xfra-layer0/terraform/$$os/terraform_$(TERRAFORM_VERSION)_$$os.zip . ; \
		unzip terraform_$(TERRAFORM_VERSION)_$$os.zip ; \
		rm terraform_$(TERRAFORM_VERSION)_$$os.zip ; \
		cd - ; \
        done

build: download-terraform
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/linux/l0-setup .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/darwin/l0-setup .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/windows/l0-setup.exe .

install-smoketest:
		for os in linux; do \
                mkdir -p build/$$os/bin ; \
                cd build/$$os/bin ; \
                aws s3 cp s3://xfra-layer0/terraform/$$os/terraform_$(TERRAFORM_VERSION)_$$os.zip . ; \
                unzip terraform_$(TERRAFORM_VERSION)_$$os.zip ; \
                rm terraform_$(TERRAFORM_VERSION)_$$os.zip ; \
                cp terraform /usr/local/bin ; \
                cd - ; \
        done
		CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/linux/l0-setup .
		cp build/linux/l0-setup ./

apply-smoketest:
		./l0-setup apply --force $(LAYER0_PREFIX)
		eval $(l0-setup endpoint -i)

destroy-smoketest:
		./l0-setup destroy --force $(LAYER0_PREFIX)

release: build

test:
	go test ./...

.PHONY: build-directories download-terraform build install-smoketest build-smoketest destroy-smoketest release test
