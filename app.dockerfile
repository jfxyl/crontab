FROM golang:1.19-alpine3.18 AS builder

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . /go/src/crontab

WORKDIR /go/src/crontab/cmd/master

RUN CGO_ENABLED=0 GOOS=linux go install -tags netgo -ldflags '-w -extldflags "-static"' ./...

WORKDIR /go/src/crontab/cmd/worker

RUN CGO_ENABLED=0 GOOS=linux go install -tags netgo -ldflags '-w -extldflags "-static"' ./...

