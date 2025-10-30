.PHONY: build clean test

GO := go
GOFLAGS := -v

build:
	mkdir -p bin
	$(GO) build $(GOFLAGS) -o bin/bsubio ./cmd/bsubio

clean:
	rm -f bsubio

test:
	$(GO) test $(GOFLAGS) ./...
