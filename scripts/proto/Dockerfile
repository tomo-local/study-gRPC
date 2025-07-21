FROM golang:1.22

# 必要なツールをインストール
RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

ENV PATH=$PATH:/go/bin

WORKDIR /workspace