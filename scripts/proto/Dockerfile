# builderステージ: Goのツールをコンパイルするためだけのステージ
FROM golang:1.22-alpine AS builder

# 必要なGoツールをインストール
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# finalステージ: 実際に利用する軽量イメージ
FROM alpine:latest

# 必要なツールをインストール
RUN apk --no-cache add protobuf protobuf-dev

# builderステージからコンパイル済みのツールをコピー
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/local/bin/

WORKDIR /workspace