# マイクロサービス連携 実行ガイド

## 🚀 セットアップと実行

### 1. 統合Docker Compose環境の起動

```bash
# プロジェクトルートで実行
cd /Users/tomo/Local/other/gRPC

# すべてのサービスを起動
docker-compose -f docker-compose-microservices.yml up -d

# ログを確認
docker-compose -f docker-compose-microservices.yml logs -f
```

### 2. サービス状態の確認

```bash
# コンテナ状態確認
docker-compose -f docker-compose-microservices.yml ps

# 各サービスのヘルスチェック
curl http://localhost:9001/health  # Auth Service (外側ポート)
curl http://localhost:9002/health  # Note Service (外側ポート)

# gRPCサービスの確認
grpcurl -plaintext localhost:9001 list      # Auth Service (外側ポート)
grpcurl -plaintext localhost:9002 list      # Note Service (外側ポート)
```

### 3. 実際の連携テスト

```bash
# Go モジュールの準備
cd /Users/tomo/Local/other/gRPC
go mod init microservices-demo
go mod tidy

# デモ実行
go run microservice_integration_demo.go
```

## 🔗 連携フロー

1. **ユーザー登録/ログイン** (Auth Service: 外側 9001 → 内側 8080)
   - ユーザー登録 → メール認証 → ログイン
   - JWTトークン取得

2. **認証付きノート作成** (Note Service: 外側 9002 → 内側 8080)
   - Authorization: Bearer {token}
   - ノート作成・取得

## 🐳 Docker Compose ネットワーク

すべてのサービスは `microservices_network` で接続:
- auth_service:8080 (内部) ← 9001:8080 (外側マッピング)
- note_service:8080 (内部) ← 9002:8080 (外側マッピング)
- auth_postgres:5432
- note_postgres:5432

## 🛠️ 開発時のTips

### サービス間通信の確認
```bash
# Auth Serviceから Note Serviceへの疎通確認
docker exec -it auth_service ping note_service

# Note Serviceから Auth Serviceへの疎通確認
docker exec -it note_service ping auth_service
```

### ログの確認
```bash
# 特定のサービスのログ
docker-compose -f docker-compose-microservices.yml logs -f auth_service
docker-compose -f docker-compose-microservices.yml logs -f note_service
```

### データベースの確認
```bash
# Auth Service DB
docker exec -it auth_postgres psql -U postgres -d auth_db

# Note Service DB
docker exec -it note_postgres psql -U noteuser -d notedb
```

## 🎯 将来的な拡張

1. **API Gateway** の追加
2. **Service Mesh** (Istio/Linkerd)
3. **分散トレーシング** (Jaeger/Zipkin)
4. **メトリクス収集** (Prometheus/Grafana)
5. **サーキットブレーカー** パターン

## 🔒 セキュリティ考慮事項

- JWTシークレットキーの共有管理
- サービス間通信のTLS暗号化
- ネットワークセグメンテーション
- 権限ベースのアクセス制御
