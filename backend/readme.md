[TOC]

# 本地安装goctl【非官方，请使用i-Things/go-zero】

1. 本地将 `go-zero 项目克隆下来：  `git clone git@github.com:i-Things/go-zero.git`
2. 到目录 `go-zero\tools\goctl 下 执行命令： `go install`
3. 后续执行下面的各种goctl命令即可

# 环境初始化

`protoc/protoc-gen-go/protoc-gen-grpc-go` 依赖可以通过下列命令 一键安装

```shell
$ goctl env check --install --verbose --force
```

# 服务新增方案

## rpc服务
```
goctl rpc new opssvr  --style=goZero -m
```
## api服务
```
goctl api new viewsvr  --style=goZero 
```

# 库表新增方案

在每个服务的 `internal/repo/relationDB` 目录下有example.go
1. 借助 `https://sql2gorm.mccode.info/` 生成对应的模型 放到 `internal/repo/relationDB/modle.go` 中
2. 复制 `internal/repo/relationDB/example.go` 到对应目录下,并修改表名
3. 将example.go中的Example替换为表名
4. 定制修改对应函数即可


# api网关接口代理模块-apisvr

```shell
#cd apisvr && goctl api go -api http/api.api  -dir ./  --style=goZero && cd ..
goctl api go -api http/backend.api  -dir ./  --style=goZero  && goctl api swagger -filename swagger.json -api http/backend.api -dir ./http&&  goctl api access  -api http/backend.api -dir ./http 
#goctl 1.9.2+
goctl api go -api http/backend.api  -dir ./  --style=goZero  && goctl api swagger -filename swagger -api http/backend.api -dir ./http&&  goctl api access  -api http/backend.api -dir ./http 

goctl api swagger -filename swagger.json -api http/backend.api -dir ./http 
goctl api access  -api http/backend.api -dir ./http 
```

# 开发环境联调

- 服务器
supos-community-edge

- 访问入口 
http://100.100.100.22:34098/home

数据库开放端口
- 100.100.100.22:34099

联调环境后端部署流程:

1. 本地build可执行文件上传到/home/supOS-V1.2.0.0-M-25102214-T5/gobackend/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o backend ./backend.go

2. 一键部署脚本(打包并推送镜像,重启backend服务)
cd /home/supOS-V1.2.0.0-M-25102214-T5/ && sh rebuild_start_backend.sh
