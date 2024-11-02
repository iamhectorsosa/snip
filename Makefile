all:
	@echo "**********************************************************"
	@echo "**                      Makefile                        **"
	@echo "**********************************************************"

build:
	CGO_ENABLED=1 go build -v -o snip .

test:
	go test ./... -cover

sec:
	gosec ./...

format:
	go fmt ./...

vet:
	go vet ./...
	staticcheck ./...
