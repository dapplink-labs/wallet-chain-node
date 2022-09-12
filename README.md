<!--
parent:
  order: false
-->

<div align="center">
  <h1> Savour Core 项目 </h1>
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

Savour HD 是 Savour 项目的钱包的 HD. 后端服务，使用 golang 编写，提供 grpc 接口给上层服务访问

**注意**: 需要 [Go 1.18+](https://golang.org/dl/)

## 安装

### 安装依赖
```bash
go mod tidy
```
### 构建程序
```bash
go build 或者 go install savour-hd
```

### 启动程序
```bash
./savour-hd -c ./config.yml
```

### 启动 RPC 接口测试界面

```bash
grpcui -plaintext 127.0.0.1:8089
```

## 贡献代码

### 第一步： fork 仓库

将 savour-core fork 到您自己的代码仓库

### 第二步： clone 您自己仓库的代码

```bash
git@github.com:guoshijiang/savour-hd.git
```

### 第三步：建立分支编写提交代码

```bash
git branch -C xxx
git checkout xxx
编写您的代码
git add .
git commit -m "xxx"
git push origin xxx
```

### 第四步：提交 PR

到你的 github 上面有一个 pr, 提交到 savour-core 代码库


### 第五步：review 完成

待 savour-core 代码维护者 review 通过之后代码会合并到 savour-core 库中，至此，您的 PR 就提交完成了 
