FROM golang:1.19-alpine3.18 AS builder

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN export GOOS=windows
RUN export GOARCH=amd64

COPY . /go/src/crontab

WORKDIR /go/src/crontab/cmd/master

RUN go install ./...

#RUN GOOS=windows GOARCH=amd64 go build -o $GOPATH/bin/master.exe -i

WORKDIR /go/src/crontab/cmd/worker

RUN go install ./...

#RUN GOOS=windows GOARCH=amd64 go build -o $GOPATH/bin/worker.exe -i