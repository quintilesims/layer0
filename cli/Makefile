SHELL:=/bin/bash
L0_VERSION?=$(shell git describe --tags)

build:
		CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/linux/l0 .
		CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/darwin/l0 .
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/windows/l0.exe .

release: build

test:
		go test ./...

install-smoketest:
		CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=$(L0_VERSION)" -a -o build/linux/l0 .
		mkdir ~/bin || :
		cp ./build/linux/l0 ~/bin
		export PATH=~/bin:$PATH

.PHONY: build release test install-smoketest
