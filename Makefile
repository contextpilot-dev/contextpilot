# ContextPilot Build System
# Cross-compile for all supported platforms

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
BINARY := contextpilot
DIST_DIR := dist

# Platforms
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all clean build release checksums

all: build

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY) .

# Build for all platforms
release: clean
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		output="$(DIST_DIR)/$(BINARY)-$$os-$$arch$$ext"; \
		echo "Building $$output..."; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o $$output . || exit 1; \
	done
	@echo "Done! Binaries in $(DIST_DIR)/"

# Create compressed archives for release
archives: release
	@cd $(DIST_DIR) && \
	for f in $(BINARY)-*; do \
		if [ -f "$$f" ]; then \
			name=$${f%.*}; \
			if echo "$$f" | grep -q ".exe"; then \
				zip "$${name}.zip" "$$f"; \
			else \
				tar -czf "$${name}.tar.gz" "$$f"; \
			fi; \
		fi; \
	done
	@echo "Archives created!"

# Generate checksums
checksums: archives
	@cd $(DIST_DIR) && shasum -a 256 *.tar.gz *.zip > checksums.txt
	@echo "Checksums generated!"

# Full release build
dist: checksums
	@echo ""
	@echo "Release artifacts:"
	@ls -la $(DIST_DIR)/

clean:
	rm -rf $(DIST_DIR)
	rm -f $(BINARY)

# Install locally
install: build
	cp $(BINARY) /usr/local/bin/

# Uninstall
uninstall:
	rm -f /usr/local/bin/$(BINARY)

# Run tests
test:
	go test -v ./...

# Development build with race detector
dev:
	go build -race -o $(BINARY) .
