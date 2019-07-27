# 生成javascript 本地签名方法

## 1. 使用go语言实现本地签名
## 2. 使用gopherjs,将go语言版本转成javascript版本
### 2.1 gopherjs安装
1. 使用go 1.12版本
2. 获取：go get -u github.com/gopherjs/gopherjs
可能会因为网络问题  \x\tools 依赖无法下载，可以从github.com/golang 中自行下载并放在gopath下。
3. 安装：go install -v github.com/gopherjs/gopherjs
4. 使用：gopherjs build WalletSendAPI.go, 在同级目录下生成对应的WalletSendAPI.js文件
