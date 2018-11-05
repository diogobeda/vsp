FROM golang:1.10.3-alpine3.8

RUN apk upgrade -U \
 && apk add bash git openssh curl \
 && rm -rf /var/cache/*

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /var/local/go
RUN mkdir bin pkg src
ENV GOPATH=/var/local/go
ENV PATH="$PATH:$GOPATH/bin"
RUN go get -u github.com/githubnemo/CompileDaemon
RUN go get -u github.com/derekparker/delve/cmd/dlv

ADD . src/github.com/diogobeda/vsp

# WORKDIR /var/local/go/src/github.com/diogobeda/vsp
# RUN dep ensure

WORKDIR /var/local/go/src/github.com/diogobeda/vsp/cmd/channel_api
CMD CompileDaemon -command="./channel_api"