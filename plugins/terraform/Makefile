build:
	go build -o terraform-provider-layer0

release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build  -a -o build/Linux/terraform-provider-layer0 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build  -a -o build/macOS/terraform-provider-layer0 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -a -o build/Windows/terraform-provider-layer0.exe .

test:
	go test ./...

.PHONY: build release test
