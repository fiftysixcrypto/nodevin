BINARY_NAME=nodevin

all: test build

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)-linux $(BINARY_NAME)-darwin $(BINARY_NAME).exe

