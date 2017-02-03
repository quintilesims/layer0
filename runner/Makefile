L0_VERSION?=$(shell git describe --tags)
CURRENT_RUNNER_DOCKER_IMAGE=quintilesims/l0-runner:$(L0_VERSION)
LATEST_RUNNER_DOCKER_IMAGE=quintilesims/l0-runner:latest

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(L0_VERSION)" -o l0-runner .
	docker build -t $(CURRENT_RUNNER_DOCKER_IMAGE) .

release: build
	docker push $(CURRENT_RUNNER_DOCKER_IMAGE)
	docker tag  $(CURRENT_RUNNER_DOCKER_IMAGE) $(LATEST_RUNNER_DOCKER_IMAGE)
	docker push $(LATEST_RUNNER_DOCKER_IMAGE)

test:
	go test ./...

.PHONY: build release test
