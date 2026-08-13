<!--
parent:
  order: false
-->

<div align="center">
  <h1> wallet-chain-node  repo </h1>
</div>

<div align="center">
  <a href="https://github.com/savour-labs/wallet-chain-node/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/savour-labs/wallet-chain-node.svg" />
  </a>
  <a href="https://github.com/savour-labs/wallet-chain-node/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/savour-labs/wallet-chain-node.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/savour-labs/wallet-chain-node">
    <img alt="GoDoc" src="https://godoc.org/github.com/savour-labs/wallet-chain-node?status.svg" />
  </a>
</div>

Savour HD is the HD of the wallet of the Savour project. The back-end service, written in golang, provides grpc interface for upper-layer service access

**Tips**: need [Go 1.18+](https://golang.org/dl/)

## Install

### Install dependencies
```bash
go mod tidy
```
### build
```bash
go build or go install wallet-chain-node
```

### start 
```bash
./wallet-chain-node -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8189
```

## Contribute

### 1.fork repo

fork wallet-chain-node to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/wallet-chain-node.git
```

### 3. create new branch and commit code

```bash
git branch -C xxx
git checkout xxx

coding

git add .
git commit -m "xxx"
git push origin xxx
```

### 4.commit PR

Have a pr on your github and submit it to the wallet-chain-node repository

### 5.review 

After the wallet-chain-node code maintainer has passed the review, the code will be merged into the wallet-chain-node library. At this point, your PR submission is complete

### 6.Disclaimer

This code has not yet been audited, and should not be used in any production systems.
