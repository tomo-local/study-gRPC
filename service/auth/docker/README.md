# Auth Service - Docker Setup

このディレクトリにはAuth serviceをDockerで実行するための設定ファイルが含まれています。

## ファイル構成

- `Dockerfile` - Auth serviceのDockerイメージを構築するためのファイル
- `docker-compose.yml` - Auth service全体（アプリ + PostgreSQL）を管理
- `.dockerignore` - Dockerビルド時に除外するファイルの設定

## 使用方法

### 1. 全体のサービスを起動

```bash
cd docker
docker-compose up -d
```

### 2. ログを確認

```bash
# 全サービスのログ
docker-compose logs -f

# Auth serviceのみのログ
docker-compose logs -f auth_service
```

### 3. サービス停止

```bash
# サービス停止
docker-compose down

# サービス停止 + データ削除
docker-compose down -v --rmi local
```

### 4. イメージの再ビルド

```bash
# 通常のビルド
docker-compose build

# キャッシュなしでビルド
docker-compose build --no-cache
```

## 環境変数の設定

環境変数は`docker-compose.yml`内で直接定義されています。必要に応じて以下を変更できます：

- **データベース設定**: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- **サーバー設定**: `SERVER_PORT`
- **JWT設定**: `JWT_SECRET_KEY`, `JWT_TOKEN_DURATION_HOURS`
- **メール設定**: `MAILER_HOST`, `MAILER_PORT`, `MAILER_USERNAME`, `MAILER_PASSWORD`

また、`../config.yml`ファイルでもデフォルト設定を管理しており、環境変数が優先されます。

## ポート

- Auth Service (gRPC): `50053`
- PostgreSQL: `5432`

## ネットワーク

サービス間は`auth_network`というDockerネットワークで通信します。将来的に他のサービスと連携する際は、同じネットワークを使用することで連携が可能です。

## 本番環境での使用

本番環境で使用する場合は、必ず以下を変更してください：

1. `docker-compose.yml`の`JWT_SECRET_KEY`を強力なものに変更
2. データベースのパスワードを変更
3. メール設定を実際の値に変更
4. 必要に応じてポートを変更

## トラブルシューティング

### サービスが起動しない場合

1. ポートが既に使用されていないか確認
2. Dockerサービスが起動しているか確認
3. ログを確認: `docker-compose logs -f`

### データベース接続エラー

1. PostgreSQLコンテナが正常に起動しているか確認
2. 環境変数が正しく設定されているか確認
3. ヘルスチェックを確認: `docker-compose ps`
