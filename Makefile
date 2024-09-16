BINARY_NAME=nodevin
MAIN_PACKAGE=./cmd/nodevin
VERSION=$(shell git describe --tags --abbrev=0)

all: test build package

build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_PACKAGE)

test:
	go test ./...

package:
	@mkdir -p release
	@cp $(BINARY_NAME)-linux release/$(BINARY_NAME)
	@tar -czvf release/$(BINARY_NAME)-linux-$(VERSION).tar.gz -C release $(BINARY_NAME)
	@shasum -a 256 release/$(BINARY_NAME)-linux-$(VERSION).tar.gz > release/$(BINARY_NAME)-linux-$(VERSION).tar.gz.sha256
	@cp $(BINARY_NAME)-macos release/$(BINARY_NAME)
	@tar -czvf release/$(BINARY_NAME)-macos-$(VERSION).tar.gz -C release $(BINARY_NAME)
	@shasum -a 256 release/$(BINARY_NAME)-macos-$(VERSION).tar.gz > release/$(BINARY_NAME)-macos-$(VERSION).tar.gz.sha256
	@cp $(BINARY_NAME).exe release/$(BINARY_NAME).exe
	@zip release/$(BINARY_NAME)-windows-$(VERSION).zip release/$(BINARY_NAME).exe
	@shasum -a 256 release/$(BINARY_NAME)-windows-$(VERSION).zip > release/$(BINARY_NAME)-windows-$(VERSION).zip.sha256
	@rm -f release/$(BINARY_NAME) release/$(BINARY_NAME).exe

clean:
	go clean
	rm -f $(BINARY_NAME)-linux $(BINARY_NAME)-macos $(BINARY_NAME).exe
	rm -rf release

.PHONY: all build test package clean
