SHELL:=/bin/bash
L0_VERSION?=$(shell git describe --tags)

build:
	go build -o l0-setup

release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/linux/l0-setup .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/darwin/l0-setup .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/windows/l0-setup.exe .

install-smoketest:
	go install
	mv $(GOPATH)/bin/setup $(GOPATH)/bin/l0-setup

test:
	go test ./...

.PHONY: build release test install-smoketest
