#!/bin/bash
set -e

# Goのbinパスを通す
export PATH="$PATH:$(go env GOPATH)/bin"

# 引数チェック
if [ $# -eq 0 ]; then
  echo "Usage: $0 <service1> [service2] ..."
  echo "Example: $0 auth memo"
  exit 1
fi

echo "Generating proto files for services: $@"

# 指定されたサービスのみ処理
for service in "$@"; do
  echo "Processing service: $service"
  
  # サービスディレクトリの存在チェック
  if [ ! -d "service/$service" ]; then
    echo "Warning: Service directory 'service/$service' not found, skipping..."
    continue
  fi
  
  # 該当サービスのprotoファイルを検索
  find "service/$service/proto" -name "*.proto" 2>/dev/null | while read proto; do
    if [ -f "$proto" ]; then
      echo "  Processing: $proto"
      dir=$(dirname "$proto")
      
      # 出力ディレクトリを設定（service/{service}/app/grpc）
      outdir="service/$service/app/grpc"
      mkdir -p "$outdir"
      
      # proto生成
      protoc \
        --proto_path="service/$service/proto" \
        --go_out="$outdir" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$outdir" \
        --go-grpc_opt=paths=source_relative \
        "$proto"
      
      echo "  Generated: $outdir"
    fi
  done || echo "  No proto files found for service: $service"
done

echo "Proto generation completed for services: $@"