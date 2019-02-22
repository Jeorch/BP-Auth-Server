FROM golang:alpine

RUN apk add --no-cache git mercurial

LABEL ddsaas.version="0.0.1" maintainer="Jeorch"

RUN go get github.com/alfredyang1986/blackmirror
RUN go get github.com/Jeorch/BP-Auth-Server

ADD deploy-config/ /go/bin/

RUN go install -v github.com/alfredyang1986/BP-Auth-Server

WORKDIR /go/bin

ENTRYPOINT ["BP-Auth-Server"]
