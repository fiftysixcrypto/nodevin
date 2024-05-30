BINARY_NAME=nodevin
MAIN_PACKAGE=./cmd/nodevin

all: test build

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_PACKAGE)

test:
	go test ./...

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_PACKAGE)

clean:
	go clean
	rm -f $(BINARY_NAME)-linux $(BINARY_NAME)-darwin $(BINARY_NAME).exe

.PHONY: all build test clean
