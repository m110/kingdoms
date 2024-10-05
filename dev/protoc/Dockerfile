FROM --platform=linux/amd64 golang:1.21

ENV PB_VERSION=25.2

RUN apt-get update && apt-get install unzip

WORKDIR /tmp
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$PB_VERSION/protoc-$PB_VERSION-linux-x86_64.zip
RUN unzip protoc-$PB_VERSION-linux-x86_64.zip -d /protobuf

ENV PATH="${PATH}:/protobuf/bin:/go/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
