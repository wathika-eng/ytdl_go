# Variables
BINARY_NAME=ytdl
GO_BUILD_FLAGS=-ldflags "-s -w"
UPX=upx

# Targets
.PHONY: all clean build-mac build-linux build-windows

all: build-mac build-linux build-windows

# Build for macOS (amd64)
build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o bin/$(BINARY_NAME)-macos ./...
	$(UPX) bin/$(BINARY_NAME)-macos

# Build for Linux (amd64)
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o bin/$(BINARY_NAME)-linux ./...
	$(UPX) bin/$(BINARY_NAME)-linux

# Build for Windows (amd64)
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o bin/$(BINARY_NAME)-windows.exe ./...
	$(UPX) bin/$(BINARY_NAME)-windows.exe

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/

run:
	@go mod tidy
	@go run cmd/main.go