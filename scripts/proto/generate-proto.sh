#!/bin/bash
set -e

# Goのbinパスを通す
export PATH="$PATH:$(go env GOPATH)/bin"

# protoファイルを探してpb.go生成
find . -name "*.proto" | while read proto; do
  dir=$(dirname "$proto")
  protoc --go_out="$dir" --go-grpc_out="$dir" --proto_path="$dir" "$proto"
done