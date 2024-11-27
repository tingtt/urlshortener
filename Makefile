GO ?= go
GOBIN ?= $(shell echo `$(GO) env GOPATH`/bin)
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)
MODULE_NAME ?= $(shell head -n1 go.mod | cut -f 2 -d ' ')

.PHONY: setup
.DEFAULT_GOAL := setup
setup:
	mkdir -p .tools
	curl -L https://github.com/dominikh/go-tools/releases/download/2024.1.1/staticcheck_darwin_arm64.tar.gz \
		-o staticcheck.tar.gz
	tar xf staticcheck.tar.gz -C .tools/
	rm staticcheck.tar.gz

.PHONY: staticcheck
staticcheck:
	.tools/staticcheck/staticcheck ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
PARALLELS ?= 10
test: fmt staticcheck
	$(GO) test ./... -short -parallel $(PARALLELS) -coverprofile=cover.out
	$(GO) tool cover -html=cover.out -o cover.html

.PHONY: testall
PARALLELS ?= 10
testall: fmt staticcheck
	$(GO) test ./... -parallel $(PARALLELS) -coverprofile=cover.out
	$(GO) tool cover -html=cover.out -o cover.html

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o urlshortener cmd/main.go
