FROM golang:1.9-alpine
# Git is needed for go get
RUN apk add --no-cache git gcc libc-dev
COPY . /go/src/github.com/svera/acquire-sackson-driver
WORKDIR /go/src/github.com/svera/acquire-sackson-driver
