SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

wallet-hd-chain:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: wallet-hd-chain

clean:
	rm wallet-hd-chain

test:
	go test -v ./...

lint:
	golangci-lint run ./...