<!--
parent:
  order: false
-->

<div align="center">
  <h1> Savour Core repo </h1>
</div>

<div align="center">
  <a href="https://github.com/SavourDao/savour-hd/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/SavourDao/savour-core.svg" />
  </a>
  <a href="https://github.com/SavourDao/savour-hd/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/SavourDao/savour-core.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/SavourDao/savour-hd">
    <img alt="GoDoc" src="https://godoc.org/github.com/SavourDao/savour-hd?status.svg" />
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
go build or go install savour-hd
```

### start 
```bash
./savour-hd -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8089
```

## Contribute

### 1.fork repo

fork savour-hd to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/savour-hd.git
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

Have a pr on your github and submit it to the savour-hd repository

### 5.review 

After the savour-hd code maintainer has passed the review, the code will be merged into the savour-hd library. At this point, your PR submission is complete
