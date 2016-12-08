L0_VERSION?=$(shell git describe --tags)
API_DOCKER_IMAGE=quintilesims/l0-api:$(L0_VERSION)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(L0_VERSION)" --tags scratch -o l0-api . 
	docker build -t $(API_DOCKER_IMAGE) .

release: build
	docker push $(API_DOCKER_IMAGE)

test:
	go test ./...
