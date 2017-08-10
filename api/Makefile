SHELL:=/bin/bash
VERSION?=$(shell git describe --tags)
CURRENT_API_DOCKER_IMAGE=quintilesims/l0-api:$(VERSION)
LATEST_API_DOCKER_IMAGE=quintilesims/l0-api:latest

deps:
	@docker login -u $(DOCKERHUB_USERNAME) -p $(DOCKERHUB_PASSWORD)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(VERSION)" -o l0-api . 
	docker build -t $(CURRENT_API_DOCKER_IMAGE) .

release: build
	docker push $(CURRENT_API_DOCKER_IMAGE)
	docker tag  $(CURRENT_API_DOCKER_IMAGE) $(LATEST_API_DOCKER_IMAGE)
	docker push $(LATEST_API_DOCKER_IMAGE)

test:
	go test ./...

.PHONY: deps build release test
