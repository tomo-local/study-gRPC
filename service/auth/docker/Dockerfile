# ビルド用の軽量なGoイメージを使用
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates

# ワーキングディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピーして依存関係をダウンロード
COPY app/go.mod app/go.sum ./
RUN go mod download

# ソースコードをコピー
COPY app/ .
COPY config.yml ./config.yml

# バイナリをビルド（静的リンク）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-server ./cmd/server/main.go

# 実行用の最小限のイメージ
FROM alpine:latest

# SSL証明書とタイムゾーンデータをインストール
RUN apk --no-cache add ca-certificates tzdata

# 非rootユーザーを作成
RUN adduser -D -s /bin/sh authuser

# ワーキングディレクトリを作成
WORKDIR /app

# ビルドしたバイナリと設定ファイルをコピー
COPY --from=builder /app/auth-server .
COPY --from=builder /app/config.yml .

# ユーザーを変更
USER authuser

# ヘルスチェック用のポート
EXPOSE 50053

# ヘルスチェックを追加
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 50053 || exit 1

# アプリケーションを実行
CMD ["./auth-server"]
