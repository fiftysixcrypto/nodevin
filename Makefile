BINARY_NAME=nodevin
MAIN_PACKAGE=./cmd/nodevin
VERSION=$(shell git describe --tags --abbrev=0)

all: build-linux-amd64 build-linux-arm64 build-macos-amd64 build-macos-arm64 build-windows-amd64 build-windows-arm64

checksum:
	shasum -a 256 $(BINARY_NAME) > $(BINARY_NAME).sha256 || true

# Build binaries for different architectures
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos-amd64 $(MAIN_PACKAGE)

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-macos-arm64 $(MAIN_PACKAGE)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o $(BINARY_NAME)-windows-arm64.exe $(MAIN_PACKAGE)

# Packaging for Linux binaries
package-linux-amd64:
	@mkdir -p release
	@mv $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	tar -czvf release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz > release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME) readme.txt

package-linux-arm64:
	@mkdir -p release
	@mv $(BINARY_NAME)-linux-arm64 $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	tar -czvf release/$(BINARY_NAME)-linux-arm64-$(VERSION).tar.gz $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-linux-arm64-$(VERSION).tar.gz > release/$(BINARY_NAME)-linux-arm64-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME) readme.txt

# Packaging for macOS binaries
package-macos-amd64:
	@mkdir -p release
	@mv $(BINARY_NAME)-macos-amd64 $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	tar -czvf release/$(BINARY_NAME)-macos-amd64-$(VERSION).tar.gz $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-macos-amd64-$(VERSION).tar.gz > release/$(BINARY_NAME)-macos-amd64-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME) readme.txt

package-macos-arm64:
	@mkdir -p release
	@mv $(BINARY_NAME)-macos-arm64 $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	tar -czvf release/$(BINARY_NAME)-macos-arm64-$(VERSION).tar.gz $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-macos-arm64-$(VERSION).tar.gz > release/$(BINARY_NAME)-macos-arm64-$(VERSION).tar.gz.sha256
	rm $(BINARY_NAME) readme.txt

# Packaging for Windows binaries
package-windows-amd64:
	@mkdir -p release
	@mv $(BINARY_NAME)-windows-amd64.exe $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	zip release/$(BINARY_NAME)-windows-amd64-$(VERSION).zip $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-windows-amd64-$(VERSION).zip > release/$(BINARY_NAME)-windows-amd64-$(VERSION).zip.sha256
	rm $(BINARY_NAME) readme.txt

package-windows-arm64:
	@mkdir -p release
	@mv $(BINARY_NAME)-windows-arm64.exe $(BINARY_NAME)
	@cp docs/cli-commands.md readme.txt
	zip release/$(BINARY_NAME)-windows-arm64-$(VERSION).zip $(BINARY_NAME) readme.txt
	shasum -a 256 release/$(BINARY_NAME)-windows-arm64-$(VERSION).zip > release/$(BINARY_NAME)-windows-arm64-$(VERSION).zip.sha256
	rm $(BINARY_NAME) readme.txt

clean:
	rm -rf release $(BINARY_NAME)-linux-* $(BINARY_NAME)-macos-* $(BINARY_NAME)-windows-*.exe $(BINARY_NAME).sha256

.PHONY: all build-linux build-macos build-windows package-linux package-macos package-windows checksum clean
