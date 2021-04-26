FROM golang:1.16-buster AS base

RUN apt-get update
RUN apt-get install git build-essential curl unzip nano -y

# Install global CLI tools
RUN go install -v golang.org/x/tools/gopls@v0.6.10
RUN go install -v github.com/go-delve/delve/cmd/dlv@v1.6.0
RUN go install -v github.com/ramya-rao-a/go-outline@v0.0.0-20200117021646-2a048b4510eb
RUN go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.37.0
RUN go install -v github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest

WORKDIR /goload

CMD gin -i -a 10000
