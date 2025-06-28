# Database Directory 📁

このディレクトリには、noteサービスのデータベース関連ファイルが格納されています。

## ファイル構成

### docker-compose.yml
- PostgreSQLコンテナの設定ファイル
- データベースの環境変数、ポート、ボリューム設定を含む

### data/
- PostgreSQLのデータ永続化用ディレクトリ
- コンテナを削除してもデータが保持される

## テーブル管理

このプロジェクトでは **GORM** を使用してテーブルの自動作成を行っています。

- `db.AutoMigrate(&model.Note{})` により、アプリケーション起動時に自動でテーブルが作成される
- 手動でのSQLマイグレーションファイルは不要
- モデル定義の変更があっても、GORMが自動で適用する

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

## データ永続化

- データは `./data` ディレクトリに保存される
- コンテナを削除しても、データは残る
- 完全にリセットしたい場合は `data` ディレクトリを削除する

## 注意事項

- `postgres_data` ボリュームにデータが永続化されます
- 初期化は初回起動時のみ実行されます
- スキーマ変更は migrations/ ディレクトリで管理してください
