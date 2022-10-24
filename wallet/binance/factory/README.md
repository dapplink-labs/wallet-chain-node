# erc20 code 

## 编译abigen二进制

```
    cd $GOPATH/src/github.com/ethereum/go-ethereum/cmd/abigen
    go build
```

## 安装solc二进制

```
    brew update
    brew upgrade
    brew tap ethereum/ethereum
    brew install cpp-ethereum
    brew linkapps cpp-ethereum
```

## 根据智能合约abi文件生成对应的合约接口文件

```
    ./abigen --sol multi_sign.sol --pkg factory --out multi_sign.go

    ./abigen --sol token.sol --pkg factory --out token.go
```