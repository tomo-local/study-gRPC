# Database Directory 📁

このディレクトリには、noteサービスのデータベース関連ファイルが格納されています。

## ファイル構成

### docker-compose.yml
- PostgreSQLコンテナの設定ファイル
- データベースの環境変数、ポート、ボリューム設定を含む

### init.sql
- Docker Composeでのコンテナ初期化時に実行されるSQLファイル
- テーブル作成、インデックス作成、サンプルデータ挿入を含む

### migrations/
- データベーススキーマのバージョン管理用ディレクトリ
- 各マイグレーションファイルは `001_`, `002_` などの連番で管理

## 使用方法

### Docker Composeでの起動
```bash
# databaseディレクトリに移動
cd database

# PostgreSQLコンテナを起動
docker-compose up -d postgres

# ログを確認
docker-compose logs postgres
```

### データベース接続確認
```bash
# PostgreSQLコンテナに接続
docker exec -it note_postgres psql -U noteuser -d notedb

# テーブル一覧を確認
\dt

# notesテーブルの内容を確認
SELECT * FROM notes;
```

## PostgreSQL設定

- **データベース名**: notedb
- **ユーザー名**: noteuser
- **パスワード**: notepass
- **ポート**: 5432

## 注意事項

- `postgres_data` ボリュームにデータが永続化されます
- 初期化は初回起動時のみ実行されます
- スキーマ変更は migrations/ ディレクトリで管理してください
