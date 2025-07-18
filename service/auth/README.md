# Auth Service

メールアドレス認証機能付きのユーザー認証マイクロサービス。

## 機能

- ユーザー登録（メールアドレス認証付き）
- ログイン/ログアウト
- JWTトークン管理
- パスワードリセット
- メールアドレス認証
- 認証確認メール再送

## 技術スタック

- Go 1.21
- gRPC
- PostgreSQL
- GORM
- JWT
- bcrypt
- Gmail SMTP

## プロジェクト構成

```
auth/
├── proto/
│   └── api/
│       └── auth.proto          # gRPC API定義
├── app/
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go         # サーバーメイン
│   │   └── client/             # クライアントサンプル
│   ├── config/
│   │   └── server/
│   │       └── config.go       # 設定管理
│   ├── db/
│   │   ├── db.go               # データベース接続
│   │   ├── query.go            # クエリ関数
│   │   └── model/
│   │       └── user.go         # データモデル
│   ├── auth/
│   │   └── auth.go             # JWT・パスワード管理
│   ├── mailer/
│   │   └── mailer.go           # メール送信
│   ├── service/
│   │   ├── auth.go             # ビジネスロジック
│   │   └── grpc_service.go     # gRPCサービス
│   └── grpc/                   # 生成されたgRPCコード
├── database/
│   ├── docker-compose.yml      # PostgreSQL
│   └── init.sql                # 初期化スクリプト
├── config.yml                  # 設定ファイル
└── README.md
```

## セットアップ

### 1. 依存関係のインストール

```bash
cd service/auth/app
go mod tidy
```

### 2. protoファイルからgRPCコードを生成

```bash
# protoファイルのディレクトリに移動
cd service/auth/proto

# gRPCコードを生成
protoc --go_out=../app/grpc --go_opt=paths=source_relative \
       --go-grpc_out=../app/grpc --go-grpc_opt=paths=source_relative \
       api/auth.proto
```

### 3. データベース起動

```bash
cd service/auth/database
docker-compose up -d
```

### 4. 設定ファイルの調整

`config.yml`を環境に合わせて調整：

```yaml
server:
  port: 50053

database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  name: auth_db
  ssl_mode: disable

jwt:
  secret_key: your-secret-key-change-this-in-production
  token_duration_hours: 24

mailer:
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: your-email@gmail.com
```

### 5. サーバー起動

```bash
cd service/auth/app
go run cmd/server/main.go
```

## 環境変数

設定ファイルの代わりに環境変数でも設定可能：

```bash
# Server
export SERVER_PORT=50053

# Database
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=auth_db
export DB_SSL_MODE=disable

# JWT
export JWT_SECRET_KEY=your-secret-key
export JWT_TOKEN_DURATION_HOURS=24

# Mailer
export MAILER_HOST=smtp.gmail.com
export MAILER_PORT=587
export MAILER_USERNAME=your-email@gmail.com
export MAILER_PASSWORD=your-app-password
export MAILER_FROM=your-email@gmail.com
```

## API

### ユーザー登録

```
rpc Register(RegisterRequest) returns (RegisterResponse);
```

### ログイン

```
rpc Login(LoginRequest) returns (LoginResponse);
```

### メールアドレス認証

```
rpc VerifyEmail(VerifyEmailRequest) returns (VerifyEmailResponse);
```

### トークン検証

```
rpc VerifyToken(VerifyTokenRequest) returns (VerifyTokenResponse);
```

### パスワードリセット要求

```
rpc RequestPasswordReset(RequestPasswordResetRequest) returns (RequestPasswordResetResponse);
```

### パスワードリセット

```
rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
```

### 認証確認メール再送

```
rpc ResendVerificationEmail(ResendVerificationEmailRequest) returns (ResendVerificationEmailResponse);
```

## 開発

### テスト

```bash
go test ./...
```

### データベース接続確認

```bash
# Adminer (http://localhost:8080)
# Server: postgres
# Username: postgres
# Password: password
# Database: auth_db
```

## セキュリティ

- パスワードはbcryptでハッシュ化
- JWTトークンは署名付き
- メール認証トークンは24時間有効
- パスワードリセットトークンは1時間有効
- SQL injection対策済み

## 本番環境での注意点

1. JWT secret keyを強力なものに変更
2. データベースのパスワードを変更
3. メール送信設定を本番環境用に変更
4. HTTPS対応
5. ログ設定の追加
6. モニタリング・アラートの設定
