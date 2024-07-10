FROM alpine:3.13
LABEL maintainer="The proxy <github.com/allenliu88>"

ARG ARCH="amd64"
ARG OS="linux"

RUN mkdir -p /opt/bin
COPY bin/nanoproxy-static-linux-amd64 /opt/bin
RUN chmod +x /opt/bin/nanoproxy-static-linux-amd64

WORKDIR /opt/bin

ENTRYPOINT [ "/opt/bin/nanoproxy-static-linux-amd64" ]