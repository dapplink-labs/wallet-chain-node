SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

savour-hd:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: savour-hd

clean:
	rm savour-hd

test:
	go test -v ./...

lint:
	golangci-lint run ./...