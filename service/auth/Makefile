.PHONY: build test clean proto docker-up docker-down run

# 変数
BINARY_NAME=auth-server
BINARY_PATH=./app/cmd/server
PROTO_PATH=./proto
GRPC_PATH=./app/grpc

# ビルド
build:
	cd app && go build -o ../bin/$(BINARY_NAME) $(BINARY_PATH)

# テスト実行
test:
	cd app && go test -v ./...

# テスト（カバレッジ付き）
test-coverage:
	cd app && go test -v -cover ./...

# クリーンアップ
clean:
	rm -rf bin/
	rm -rf app/grpc/api/

# protoファイルからgRPCコードを生成
proto:
	cd $(PROTO_PATH) && protoc --go_out=../app/grpc --go_opt=paths=source_relative \
		--go-grpc_out=../app/grpc --go-grpc_opt=paths=source_relative \
		api/auth.proto

# Dockerコンテナを起動
docker-up:
	cd database && docker-compose up -d

# Dockerコンテナを停止
docker-down:
	cd database && docker-compose down

# Dockerコンテナを停止してボリュームも削除
docker-down-clean:
	cd database && docker-compose down -v

# サーバーを起動
run: proto build
	./bin/$(BINARY_NAME)

# 開発用サーバーを起動
dev: proto
	cd app && go run $(BINARY_PATH)

# 依存関係の更新
deps:
	cd app && go mod tidy

# 依存関係のダウンロード
deps-download:
	cd app && go mod download

# フォーマット
fmt:
	cd app && go fmt ./...

# lintチェック
lint:
	cd app && golangci-lint run

# 全体のセットアップ
setup: deps proto docker-up
	@echo "Auth service setup complete!"

# ヘルプ
help:
	@echo "Available targets:"
	@echo "  build         - Build the auth server binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  proto         - Generate gRPC code from proto files"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  run           - Build and run the server"
	@echo "  dev           - Run the server in development mode"
	@echo "  deps          - Update dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  setup         - Complete setup for development"
	@echo "  help          - Show this help message"
