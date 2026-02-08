BINARY = bilingual_pdf
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: build test test-integration test-cover lint clean install

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -f testdata/*.pdf testdata/*.html
