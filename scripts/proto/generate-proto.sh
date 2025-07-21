#!/bin/bash
set -e

for service in "$@"; do
  find service/${service}/proto/**/*.proto | while read proto; do
    echo "Processing proto file: $proto"
    echo "Service: ${service}"

    outdir="service/${service}/app/grpc"
    mkdir -p "$outdir"

    # protoファイルのディレクトリをモジュールルートとして使用
    protodir="service/${service}/proto"

    echo "Generating Go code for $proto in $outdir"
    protoc \
      --proto_path="$protodir" \
      --go_out="$outdir" \
      --go_opt=paths=source_relative \
      --go-grpc_out="$outdir" \
      --go-grpc_opt=paths=source_relative \
      "$proto"
  done
done