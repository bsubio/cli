.PHONY: build build-static clean test release

GO := go
GOFLAGS := -v
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

build:
	mkdir -p bin
	$(GO) build $(GOFLAGS) -o bin/bsubio ./cmd/bsubio

build-static:
	mkdir -p bin
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o bin/bsubio ./cmd/bsubio

release:
	@echo "Building static binaries for release..."
	@mkdir -p bin/release
	@for os in linux darwin; do \
		for arch in amd64 arm64; do \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 $(GO) build -ldflags="$(LDFLAGS)" -o bin/release/bsubio-$$os-$$arch ./cmd/bsubio; \
		done; \
	done
	@echo "Building windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags="$(LDFLAGS)" -o bin/release/bsubio-windows-amd64.exe ./cmd/bsubio
	@echo "Release binaries built in bin/release/"

clean:
	rm -rf bin
	rm -f bsubio

test:
	$(GO) test $(GOFLAGS) ./...
