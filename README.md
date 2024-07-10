# nanoproxy-static-static

This is a tiny HTTP forward proxy written in Go, for me to gain experience in the Go language.

This proxy accepts all requests and forwards them directly to the origin/target server. It performs no caching.

Despite this not being a full proxy implementation, it is blazing fast. In particular it is significantly faster than Squid and slightly faster than Apache's mod_proxy. This demonstrates that Go's built-in HTTP library is of a very high quality and that the Go runtime is quite performant.

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
nohup bin/nanoproxy-static-darwin-arm64 --port 9100 --target example.com:9100 > nanoproxy-static.log 2>&1 &
```

For Linux

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/nanoproxy-static-linux-amd64

chmod +x bin/nanoproxy-static-linux-amd64
nohup bin/nanoproxy-static-linux-amd64 --port 9100 --target example.com:9100 > nanoproxy-static.log 2>&1 &
```

For Windows

```shell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/nanoproxy-static-windows-amd64.exe

chmod +x bin/nanoproxy-static-windows-amd64.exe
nohup bin/nanoproxy-static-windows-amd64.exe --port 9100 --target example.com:9100 > nanoproxy-static.log 2>&1 &
```

Or just run

```shell
go run main.go --port 9100 --target example.com:9100
```

Validate port 9100 is open and open the firewall

```shell
## 查看运行详情，默认是监听8080端口，可通过--port参数指定，例如，本例中指定为9100
netstat -tunlp | grep 9100
tcp6       0      0 :::9100                 :::*                    LISTEN      19598/./nanoproxy-static   

## 开放端口
firewall-cmd --zone=public --add-port=9100/tcp --permanent
## 重载配置
firewall-cmd --reload
```

## Validation

请求转发验证：

```shell
curl 1.1.1.1:9100/metrics
```

其中，`1.1.1.1`是`nanoproxy-static`所在节点，如上请求将会被转发至`example.com:9100/metrics`。
