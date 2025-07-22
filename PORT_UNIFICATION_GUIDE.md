# ポート統一アーキテクチャ説明

## 🎯 ポート統一のメリット

### 従来の方式 (Before)
```
localhost:8080 → auth_service:8080
localhost:8082 → note_service:8082
localhost:8084 → greeting_service:8084
```

### 統一後の方式 (After)
```
localhost:9001 → auth_service:8080
localhost:9002 → note_service:8080
localhost:9003 → greeting_service:8080  (将来追加時)
```

## ✅ この方式の利点

### 1. **統一性** 🎨
- 全サービスが内部で同じポート8080を使用
- 設定ファイルやDockerfileの統一化
- 運用・保守が簡単

### 2. **スケーラビリティ** 📈
- 新しいサービス追加時は外側ポートを増やすだけ
- 内部設定は全て同じでOK
- Kubernetesへの移行も容易

### 3. **開発効率** ⚡
- テンプレート化しやすい
- 設定の間違いが減る
- 新人エンジニアにも分かりやすい

## 🔧 実際の設定

### docker-compose.yml
```yaml
services:
  auth_service:
    environment:
      SERVER_PORT: 8080  # 内部ポート統一
    ports:
      - "9001:8080"      # 外側:内側マッピング

  note_service:
    environment:
      SERVER_PORT: 8080  # 内部ポート統一
    ports:
      - "9002:8080"      # 外側:内側マッピング
```

### サービス間通信
```go
// コンテナ内からの通信（Docker network経由）
authConn := grpc.NewClient("auth_service:8080")  // 内部ポート

// 外部からの通信（開発・デバッグ用）
authConn := grpc.NewClient("localhost:9001")     // 外側ポート
```

## 🌐 ネットワーク構成

```
External Access:
localhost:9001 ──┐
                 ├─► Docker Network (microservices_network)
localhost:9002 ──┘    │
                      ├─ auth_service:8080
                      ├─ note_service:8080
                      ├─ auth_postgres:5432
                      └─ note_postgres:5432

Internal Communication:
auth_service:8080 ←→ note_service:8080 (直接通信)
```

## 🚀 運用での活用例

### Load Balancer設定
```nginx
upstream auth_cluster {
    server localhost:9001;
    server localhost:9004;  # auth_service_2
    server localhost:9007;  # auth_service_3
}

upstream note_cluster {
    server localhost:9002;
    server localhost:9005;  # note_service_2
    server localhost:9008;  # note_service_3
}
```

### Kubernetes移行時
```yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-service
spec:
  ports:
  - port: 8080        # 統一ポート
    targetPort: 8080  # 統一ポート
  selector:
    app: auth-service
```

## 💡 ベストプラクティス

1. **9001番台**: 認証系サービス
2. **9002番台**: データ管理系サービス
3. **9003番台**: 通知・メール系サービス
4. **9004番台**: ファイル・アップロード系

この統一により、より保守性・拡張性の高いマイクロサービス アーキテクチャが実現できます！
