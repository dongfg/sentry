FROM golang:1.18-alpine AS BUILDER

WORKDIR /go/src/app
ENV GOPROXY=https://goproxy.cn
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && set -eux \
    && apk --no-cache add build-base
COPY . .
RUN go mod download
RUN go build -a -o sentry

FROM alpine:latest
LABEL MAINTAINER "dongfg <mail@dongfg.com>"

ENV USER_NAME=runner
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && set -eux \
    && adduser -s /sbin/nologin -u 1000 -D ${USER_NAME} \
    && apk --no-cache add ca-certificates tzdata mysql-client mariadb-connector-c nano dumb-init \
    && mkdir -p /etc/sentry/

COPY --from=BUILDER /go/src/app/sentry /usr/local/bin
USER ${USER_NAME}

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/usr/local/bin/sentry", "app"]
