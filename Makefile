BINARY_NAME=nodevin
MAIN_PACKAGE=./cmd/nodevin
VERSION=$(shell git describe --tags --abbrev=0)

all: checksum build-linux build-macos build-windows package-linux package-macos package-windows

checksum:
	shasum -a 256 $(BINARY_NAME) > $(BINARY_NAME).sha256 || true

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos $(MAIN_PACKAGE)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows.exe $(MAIN_PACKAGE)

package-linux:
	@mkdir -p release
	@mv $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)
	tar -czvf release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz $(BINARY_NAME)
	shasum -a 256 release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz > release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME)

package-macos:
	@mkdir -p release
	@mv $(BINARY_NAME)-macos $(BINARY_NAME)
	tar -czvf release/$(BINARY_NAME)-macos-$(VERSION).tar.gz $(BINARY_NAME)
	shasum -a 256 release/$(BINARY_NAME)-macos-$(VERSION).tar.gz > release/$(BINARY_NAME)-macos-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME)

package-windows:
	@mkdir -p release
	@mv $(BINARY_NAME)-windows.exe $(BINARY_NAME).exe
	zip release/$(BINARY_NAME)-windows-$(VERSION).zip $(BINARY_NAME).exe
	shasum -a 256 release/$(BINARY_NAME)-windows-$(VERSION).zip > release/$(BINARY_NAME)-windows-$(VERSION).zip.sha256
	rm $(BINARY_NAME).exe

clean:
	rm -rf release $(BINARY_NAME)-linux $(BINARY_NAME)-macos $(BINARY_NAME)-windows.exe $(BINARY_NAME).sha256

.PHONY: all build-linux build-macos build-windows package-linux package-macos package-windows checksum clean
