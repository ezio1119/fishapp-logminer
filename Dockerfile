FROM golang:1.13-alpine AS dev

WORKDIR /app
RUN apk add --no-cache alpine-sdk && \
    go get github.com/pilu/fresh \
          golang.org/x/tools/gopls \
          github.com/mdempsky/gocode \
          github.com/uudashr/gopkgs/v2/cmd/gopkgs \
          github.com/ramya-rao-a/go-outline \
          github.com/stamblerre/gocode \
          github.com/rogpeppe/godef \
          github.com/sqs/goreturns \ 
          golang.org/x/lint/golint \
          github.com/go-delve/delve/cmd/dlv