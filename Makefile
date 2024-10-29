all:
	@echo "**********************************************************"
	@echo "**                      Makefile                        **"
	@echo "**********************************************************"

build:
	go build -o snippets

test:
	go test ./... -cover

format:
	go fmt ./...
