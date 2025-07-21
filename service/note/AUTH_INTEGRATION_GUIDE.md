# Note Service - Auth Service 連携設定ガイド

## 🔧 環境変数設定

### Note Service 用の Auth Service 連携設定

```bash
# Auth Service 接続設定
export AUTH_SERVICE_HOST=auth_service        # Auth Service のホスト名
export AUTH_SERVICE_PORT=8080                # Auth Service のポート（内部）
export AUTH_SERVICE_URL=auth_service:8080    # 完全なURL（優先）
export AUTH_ENABLED=true                     # 認証機能の有効/無効
export AUTH_TIMEOUT_SECONDS=10               # タイムアウト時間（秒）

# JWT設定（Auth Serviceと共有）
export JWT_SECRET_KEY=your-shared-secret-key

# Note Service 自体の設定
export SERVER_PORT=8080                      # Note Service のポート
export DB_HOST=note_postgres                 # データベースホスト
export DB_PORT=5432                          # データベースポート
export DB_USER=noteuser                      # データベースユーザー
export DB_PASSWORD=notepass                  # データベースパスワード
export DB_NAME=notedb                        # データベース名
export DB_SSL_MODE=disable                   # SSL モード
```

## 🐳 Docker Compose での設定例

```yaml
services:
  note_service:
    environment:
      # Note Service 基本設定
      SERVER_PORT: 8080
      DB_HOST: note_postgres
      DB_PORT: 5432
      DB_USER: noteuser
      DB_PASSWORD: notepass
      DB_NAME: notedb
      DB_SSL_MODE: disable

      # Auth Service 連携設定
      AUTH_SERVICE_HOST: auth_service
      AUTH_SERVICE_PORT: 8080
      AUTH_SERVICE_URL: auth_service:8080
      AUTH_ENABLED: true
      AUTH_TIMEOUT_SECONDS: 10
      JWT_SECRET_KEY: microservices-shared-secret-key
    ports:
      - "9002:8080"
    depends_on:
      - auth_service
```

## 🌐 接続パターン

### 1. Docker Compose 環境
```go
// コンテナ間通信（内部ネットワーク）
authClient := grpc.NewClient("auth_service:8080")
```

### 2. ローカル開発環境
```go
// 外部からアクセス（ポートマッピング経由）
authClient := grpc.NewClient("localhost:9001")
```

### 3. 環境変数を使用した柔軟な接続
```go
// 環境変数から取得
authServiceURL := os.Getenv("AUTH_SERVICE_URL")
if authServiceURL == "" {
    authServiceURL = "localhost:9001"  // デフォルト
}
authClient := grpc.NewClient(authServiceURL)
```

## 🔐 認証フローの設定

### Note Service での JWT 検証
```go
// 1. クライアントから JWT トークンを受信
// 2. Auth Service に トークン検証を依頼
// 3. 有効な場合のみ Note の操作を許可

// 設定例
type AuthConfig struct {
    ServiceURL     string  `env:"AUTH_SERVICE_URL"`
    Enabled        bool    `env:"AUTH_ENABLED"`
    TimeoutSeconds int     `env:"AUTH_TIMEOUT_SECONDS"`
    JWTSecretKey   string  `env:"JWT_SECRET_KEY"`
}
```

## 🎯 使用例

### クライアント側（統合デモ）
```bash
# 環境変数で接続先を指定
export AUTH_SERVICE_URL=localhost:9001
export NOTE_SERVICE_URL=localhost:9002

# デモ実行
go run example_auth_note_integration.go
```

### サーバー側（Note Service）
```bash
# Auth Service連携を有効化
export AUTH_ENABLED=true
export AUTH_SERVICE_URL=auth_service:8080
export JWT_SECRET_KEY=shared-secret-key

# Note Service 起動
go run cmd/server/main.go
```

## 🛠️ トラブルシューティング

### Auth Service に接続できない場合
```bash
# 1. Auth Service が起動しているか確認
docker-compose ps auth_service

# 2. ネットワーク接続確認
docker exec note_service ping auth_service

# 3. Auth Service のヘルスチェック
curl http://localhost:9001/health
```

### 認証が失敗する場合
```bash
# 1. JWT Secret Key が一致しているか確認
echo $JWT_SECRET_KEY

# 2. トークンの有効期限確認
# 3. Auth Service のログ確認
docker-compose logs auth_service
```

## 📈 スケーリング時の考慮事項

### Load Balancer 使用時
```yaml
# 複数の Auth Service インスタンス
AUTH_SERVICE_URL: "auth-lb.example.com:443"
AUTH_TIMEOUT_SECONDS: 5  # LB経由では短めに設定
```

### Service Mesh 使用時
```yaml
# サービスメッシュ経由での通信
AUTH_SERVICE_URL: "auth-service.default.svc.cluster.local:8080"
```

この設定により、Note Service が柔軟に Auth Service と連携できるようになります！🚀
