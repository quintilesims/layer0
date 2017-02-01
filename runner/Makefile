SHELL:=/bin/bash
L0_VERSION?=$(shell git describe --tags)
RUNNER_DOCKER_IMAGE=quintilesims/l0-runner:$(L0_VERSION)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(L0_VERSION)" -o l0-runner .
	docker build -t $(RUNNER_DOCKER_IMAGE) .

release: build
	docker push $(RUNNER_DOCKER_IMAGE)

test:
	go test ./...

.PHONY: build release test
