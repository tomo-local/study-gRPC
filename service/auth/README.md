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

**⚠️ 重要**: データベースが起動していないと以下のようなエラーが発生します：
```
failed to connect to `user=postgres database=auth_db`: dial tcp [::1]:5432: connect: connection refused
```

この場合は以下のコマンドでデータベースの状態を確認してください：

```bash
# データベースコンテナの状態確認
docker-compose ps

# ログを確認
docker-compose logs postgres

# データベースを再起動
docker-compose restart postgres

# 完全に停止して再起動
docker-compose down && docker-compose up -d
```

### 4. 環境変数の設定

このプロジェクトは環境変数で設定を管理します。開発環境でのセットアップ手順：

#### 4-1. 環境変数ファイルをコピー

**最初に必ずこの手順を実行してください！**

```bash
# service/auth ディレクトリに移動
cd service/auth

# .env.example を .env にコピー
cp .env.example .env
```

**注意**: この手順を実行しないと、必須の環境変数が設定されずにアプリケーションが起動しません。

#### 4-2. 必要に応じて設定値を変更

`.env`ファイルを開いて、あなたの環境に合わせて値を変更してください：

```bash
# 特に以下の項目は必ず変更してください 🔒
JWT_SECRET_KEY=your-super-secret-key-change-this-in-production
DB_PASSWORD=your-database-password
MAILER_USERNAME=your-email@gmail.com
MAILER_PASSWORD=your-app-password
MAILER_FROM=your-email@gmail.com
```

#### 4-3. セキュリティ注意事項 ⚠️

- `.env`ファイルは`.gitignore`に含まれているため、Gitにコミットされません
- 本番環境では必ず強固な秘密鍵とパスワードを使用してください
- メールのパスワードにはアプリパスワードを使用することをおすすめします

### 5. サーバー起動

```bash
cd service/auth/app
go run cmd/server/main.go
```

## 環境変数

このプロジェクトは環境変数で設定を管理します。

### 必須環境変数 📋

以下の環境変数は必須です（設定されていないとアプリが起動しません）：

- `DB_HOST` - データベースホスト
- `DB_USER` - データベースユーザー
- `DB_PASSWORD` - データベースパスワード
- `DB_NAME` - データベース名
- `JWT_SECRET_KEY` - JWT署名用の秘密鍵
- `MAILER_HOST` - メールサーバーホスト
- `MAILER_PORT` - メールサーバーポート
- `MAILER_USERNAME` - メールアカウントのユーザー名
- `MAILER_PASSWORD` - メールアカウントのパスワード
- `MAILER_FROM` - 送信元メールアドレス

### 全環境変数一覧

設定ファイルの代わりに環境変数でも設定可能：

```bash
# Server
export SERVER_PORT=8080

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

### よくあるトラブルシューティング 🔧

#### データベース接続エラー
```
failed to connect to `user=postgres database=auth_db`: dial tcp [::1]:5432: connect: connection refused
```

**解決方法:**
1. PostgreSQLコンテナが起動しているか確認
   ```bash
   cd service/auth/database
   docker-compose ps
   ```

2. コンテナを再起動
   ```bash
   docker-compose restart postgres
   ```

3. ポート5432が他のプロセスで使用されていないか確認
   ```bash
   lsof -i :5432
   ```

#### 環境変数設定エラー
```
failed to process environment variables: required key [環境変数名] missing value
```

**解決方法:**
1. `.env`ファイルが存在することを確認
   ```bash
   ls -la .env
   ```

2. `.env.example`から`.env`をコピー
   ```bash
   cp .env.example .env
   ```

3. 必須の環境変数が設定されているか確認

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
