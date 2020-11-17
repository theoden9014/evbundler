GOARCH		?= $(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))
GOOS		?= $(shell uname | tr A-Z a-z)
GO		?= GOOS=$(GOOS) GOARCH=$(GOARCH) go

PKGS		:= $(shell go list -f '{{.Dir}}' ./...)
SOURCES	:= $(foreach dir, $(PKGS), $(wildcard $(dir)/*.go))

.DEFAULT_GOAL := build

.PHONY: test
test:
	$(GO) test -race -v $(PKGS)

.PHONY: build
build: bin/evbundler

bin/evbundler: $(SOURCES)
	CGO_ENABLED=0 $(GO) build -o ./$@ ./cmd/evbundler

clean:
	@rm -f bin/*
