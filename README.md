# nanoproxy-static

This is a tiny HTTP forward proxy written in Go, which accepts all requests and forwards them directly to the origin/target server host which is set by `--target` flag. It performs no caching.

> Despite this not being a full proxy implementation, it is blazing fast. In particular it is significantly faster than Squid and slightly faster than Apache's mod_proxy. This demonstrates that Go's built-in HTTP library is of a very high quality and that the Go runtime is quite performant.

## Prerequisites

- go 1.22.3
- toolchain go1.22.5

## Build & Run

Clone

```shell
git clone https://github.com/allenliu88/nanoproxy-static.git
cd nanoproxy-static
```

For Mac Apple M1

```shell
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/nanoproxy-static-darwin-arm64

chmod +x bin/nanoproxy-static-darwin-arm64
nohup bin/nanoproxy-static-darwin-arm64 --port 9100 --target example.com:9100 --verbose > nanoproxy-static.log 2>&1 &
```

For Linux

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/nanoproxy-static-linux-amd64

chmod +x bin/nanoproxy-static-linux-amd64
nohup bin/nanoproxy-static-linux-amd64 --port 9100 --target example.com:9100 --verbose > nanoproxy-static.log 2>&1 &
```

For Windows

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/nanoproxy-static-windows-amd64.exe

chmod +x bin/nanoproxy-static-windows-amd64.exe
nohup bin/nanoproxy-static-windows-amd64.exe --port 9100 --target example.com:9100 --verbose > nanoproxy-static.log 2>&1 &
```

Or just run

```shell
go run main.go --port 9100 --target example.com:9100 --verbose
```

Validate port 9100 is open and open the firewall

```shell
## 查看运行详情，默认是监听8080端口，可通过--port参数指定，例如，本例中指定为9100
$ netstat -tunlp | grep 9100
tcp6       0      0 :::9100                 :::*                    LISTEN      19598/./nanoproxy-static   

## 注意，`+c 50`是指COMMAND列宽度
$ lsof -i :9100 +c 50
COMMAND                         PID  USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
nanoproxy-static-darwin-arm64 13293 allen    5u  IPv6 0x730791ba7d146db1      0t0  TCP *:hp-pdl-datastr (LISTEN)

## 开放端口
$ firewall-cmd --zone=public --add-port=9100/tcp --permanent
## 重载配置
$ firewall-cmd --reload

## 杀掉进程
$ kill -9 13293
```

## Validation

请求转发验证：

```shell
curl localhost:9100/metrics
```

其中，`localhost`是`nanoproxy-static`所在节点，如上请求将会被转发至`example.com:9100/metrics`。

## Kubernetes

### Docker

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/nanoproxy-static-linux-amd64
docker image build --platform linux/amd64 -t allen88/nanoproxy-static:0.1.0 -f ./Dockerfile .
docker push allen88/nanoproxy-static:0.1.0
docker tag allen88/nanoproxy-static:0.1.0 harbor.open.hand-china.com/hskp/nanoproxy-static:0.1.0
docker push harbor.open.hand-china.com/hskp/nanoproxy-static:0.1.0
```

### Deployment

```shell
kubectl apply -f ./Deployment.yml -n istio-test
```

## Ansible

```shell
cd bin
tar czvf nanoproxy-static-linux-amd64-0.1.0.tar.gz nanoproxy-static-linux-amd64 
```

## Use Case

### 主机监控

原始地址：`example.com:9100/metrics`
启动命令：`go run main.go --port 9100 --target example.com:9100 --verbose`
验证命令：`curl localhost:9100/metrics`


### 服务监控

原始地址：`example.com:15020/stats/prometheus`
端口转换：`15020 -> 30654(NodePort)`
暴露方式：`istio/opentelemetry-javaagent`
启动命令：`go run main.go --port 15020 --target example.com:30654 --verbose`
验证命令：`curl localhost:15020/stats/prometheus`
