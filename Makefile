BINARY := gl
VERSION := 2.0.0

.PHONY: build install clean test

build:
	go build -o $(BINARY) ./cmd/gl/

install:
	go install ./cmd/gl/

clean:
	rm -f $(BINARY)

test:
	go test ./...
