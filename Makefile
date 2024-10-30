all:
	@echo "**********************************************************"
	@echo "**                      Makefile                        **"
	@echo "**********************************************************"

build:
	CGO_ENABLED=1 go build -v -o snippets .

test:
	go test ./... -cover

format:
	go fmt ./...
