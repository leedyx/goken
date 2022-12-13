FROM golang:1.19.4-alpine3.17 as builder
MAINTAINER leedyx1990

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /opt/go/goken
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o goken .

FROM alpine  as final

# 时区设置成当前时区
#RUN apk add --no-cache tzdata
ENV TZ="Asia/Shanghai"
# 移动到用于存放生成的二进制文件的 /opt/app 目录
WORKDIR /opt/app
# 将二进制文件从 /opt/gat1400-Go/api-server 目录复制到这里
COPY --from=builder /opt/go/goken/goken .
# 在容器目录 /opt/app 创建一个目录 为config
#RUN mkdir config .
#COPY --from=builder  /opt/gat1400-Go/api-server/config/app.json ./config/

# 指定运行时环境变量
ENV GIN_MODE=release \
    PORT=38080
VOLUME ["/opt/app/log","/opt/app/token"]
# 声明服务端口
EXPOSE 38080
ENTRYPOINT ["/opt/app/goken"]