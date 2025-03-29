# gRPC 学習リポジトリ

このリポジトリは、以下のサイトで紹介されている gRPC の学習内容を実践するためのものです：

[【gRPC入門】Protocol BuffersとgRPCを使ったAPI開発](https://note.com/shunex/n/nd8109a1144a5)

## 概要

gRPC は Google が開発した高性能な RPC フレームワークで、Protocol Buffers（protobuf）を使用して効率的な通信を実現します。このリポジトリでは、以下の内容を学習・実践します：

- Protocol Buffers を使ったデータ定義
- gRPC サーバーとクライアントの実装
- gRPC を使った双方向ストリーミング通信

## セットアップ

1. 必要なツールをインストールします：
   - `protoc`（Protocol Buffers コンパイラ）
   - `protoc-gen-go` と `protoc-gen-go-grpc`（Go 用のコード生成プラグイン）

   ```bash
   brew install protobuf
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

2. プロジェクトの依存関係をインストールします：

   ```bash
   go mod tidy
   ```

3. Protocol Buffers ファイルからコードを生成します：

   ```bash
   protoc --go_out=. --go-grpc_out=. chat/chat.proto
   ```

## 実行方法

1. gRPC サーバーを起動します：

   ```bash
   go run server/main.go
   ```

2. gRPC クライアントを実行します：

   ```bash
   go run client/main.go
   ```

## 参考リンク

- [gRPC 公式ドキュメント](https://grpc.io/docs/)
- [Protocol Buffers 公式ドキュメント](https://protobuf.dev/)

---

このリポジトリを通じて、gRPC の基礎から応用までしっかり学んでいこう！💪✨
