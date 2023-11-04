SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

wallet-chain-node:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: wallet-chain-node

clean:
	rm wallet-chain-node

test:
	go test -v ./...

lint:
	golangci-lint run ./...