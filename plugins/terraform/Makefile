build:
	go build -o terraform-provider-layer0

release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build  -a -o build/linux/terraform-provider-layer0 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build  -a -o build/darwin/terraform-provider-layer0 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -a -o build/windows/terraform-provider-layer0.exe .

.PHONY: build release
