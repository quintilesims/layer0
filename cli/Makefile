SHELL:=/bin/bash
L0_VERSION?=$(shell git describe --tags)

build: 
	go build -o l0

release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/Linux/l0 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/macOS/l0 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/Windows/l0.exe .

test:
	go test ./...

install-smoketest:
	go install
	mv $(GOPATH)/bin/cli $(GOPATH)/bin/l0

.PHONY: build release test install-smoketest
