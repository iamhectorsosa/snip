all:
	@echo "**********************************************************"
	@echo "**                      Makefile                        **"
	@echo "**********************************************************"

run:
	go run .

test:
	go test ./... -cover

format:
	go fmt ./...
