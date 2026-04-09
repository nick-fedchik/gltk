BINARY := gl
VERSION := 2.0.0

.PHONY: build install clean test

build:
	$(MAKE) build-linux
	$(MAKE) build-windows

build-linux:
	set GOOS=linux&& set GOARCH=amd64&& go build -o $(BINARY) ./cmd/gl/

build-windows:
	set GOOS=windows&& set GOARCH=amd64&& go build -o $(BINARY).exe ./cmd/gl/

install:
	go install ./cmd/gl/

clean:
	rm -f $(BINARY)

test:
	go test ./...
