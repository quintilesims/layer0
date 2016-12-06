L0_VERSION?=$(shell git describe --tags)
TERRAFORM_VERSION=0.7.13

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

release: build

test:
	go test ./...

.PHONY: build-directories download-terraform build release test
