<!--
parent:
  order: false
-->

<div align="center">
  <h1> hdwalet-chain  repo </h1>
</div>

<div align="center">
  <a href="https://github.com/savour-labs/wallet-hd-chain/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/savour labs/savour-core.svg" />
  </a>
  <a href="https://github.com/savour-labs/wallet-hd-chain/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/savour labs/savour-core.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/savour-labs/wallet-hd-chain">
    <img alt="GoDoc" src="https://godoc.org/github.com/savour-labs/wallet-hd-chain?status.svg" />
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
go build or go install wallet-hd-chain
```

### start 
```bash
./wallet-hd-chain -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8089
```

## Contribute

### 1.fork repo

fork wallet-hd-chain to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/wallet-hd-chain.git
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

Have a pr on your github and submit it to the wallet-hd-chain repository

### 5.review 

After the wallet-hd-chain code maintainer has passed the review, the code will be merged into the wallet-hd-chain library. At this point, your PR submission is complete

### 6.Disclaimer

This code has not yet been audited, and should not be used in any production systems.
