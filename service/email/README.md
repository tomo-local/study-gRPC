# Email Service

学習用メールサーバーサービス

## 概要

Go言語で実装されたSMTP/IMAPメールサーバーです。学習目的で作成されており、基本的なメール送受信機能を提供します。

## 技術スタック

- **言語**: Go
- **SMTP実装**: `github.com/emersion/go-smtp`
- **IMAP実装**: `github.com/emersion/go-imap` (Phase 2)
- **データベース**: PostgreSQL
- **ORM**: `gorm.io/gorm`
- **コンテナ**: Docker & Docker Compose

## アーキテクチャ

```
┌─────────────────┐    SMTP     ┌──────────────────┐
│   Webアプリ     │ ─────────→ │  Go SMTP Server  │
└─────────────────┘             └──────────────────┘
                                          │
                                          ▼
                                ┌──────────────────┐
                                │   PostgreSQL     │
                                │  (認証・ユーザー) │
                                └──────────────────┘
```

## 実装フェーズ

### Phase 1: SMTP送信サーバー
- [x] 基本的なSMTP機能
- [x] PostgreSQL認証連携
- [x] 既存アプリからのメール送信

### Phase 2: IMAP受信サーバー（オプション）
- [ ] メール受信・保存機能
- [ ] メールボックス管理

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.19+

### 起動方法

```bash
# データベース起動
cd database
docker-compose up -d

# サーバー起動
cd app
go run cmd/server/main.go
```

### テスト用クライアント

```bash
cd app
go run cmd/client/main.go
```

## 設定

### 環境変数

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=email_user
DB_PASSWORD=email_pass
DB_NAME=email_db

# SMTP Server
SMTP_HOST=localhost
SMTP_PORT=2525
SMTP_TLS=false

# IMAP Server (Phase 2)
IMAP_HOST=localhost
IMAP_PORT=1143
IMAP_TLS=false
```

## API仕様

### SMTP機能
- メール送信
- ユーザー認証
- メールキューイング

### IMAP機能 (Phase 2)
- メール受信
- メールボックス管理
- メール検索

## 注意事項

⚠️ **重要**: このメールサーバーは学習目的で作成されており、本番環境での使用は想定していません。

- セキュリティ機能は最小限
- ローカル開発環境でのみ動作テスト済み
- 実際のメール配信には使用しないでください

## 参考資料

- [RFC 5321 - SMTP](https://tools.ietf.org/html/rfc5321)
- [RFC 3501 - IMAP](https://tools.ietf.org/html/rfc3501)
- [go-smtp Documentation](https://github.com/emersion/go-smtp)
- [go-imap Documentation](https://github.com/emersion/go-imap)
