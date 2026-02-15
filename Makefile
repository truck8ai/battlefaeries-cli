BINARY_NAME=bf
VERSION?=1.0.0
BUILD_DIR=dist
MODULE=github.com/truck8ai/battlefaeries-cli
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: build build-all install clean

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

build-all: clean
	mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

clean:
	rm -rf $(BUILD_DIR) $(BINARY_NAME)
