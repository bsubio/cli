.PHONY: build clean test

GO := go
GOFLAGS := -v

build:
	$(GO) build $(GOFLAGS) -o bsubio ./cmd/bsubio

clean:
	rm -f bsubio

test:
	$(GO) test $(GOFLAGS) ./...
