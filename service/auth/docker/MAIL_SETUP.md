# メール送信設定ガイド 📧

このガイドでは、Auth ServiceでGmail SMTPを使用して実際にメールを送信する設定方法を説明します。

## 1. Googleアカウントの準備

### ステップ1: 2段階認証を有効にする
1. [Google アカウント](https://myaccount.google.com/)にアクセス
2. 「セキュリティ」タブをクリック
3. 「2段階認証」を有効にする

### ステップ2: アプリパスワードを生成する
1. Googleアカウントの「セキュリティ」で「アプリパスワード」をクリック
2. アプリを選択: 「メール」
3. デバイスを選択: 「その他（カスタム名）」→「gRPC Auth Service」
4. 生成された16桁のパスワードをメモ（例: `abcd efgh ijkl mnop`）

## 2. 環境変数の設定

### `.env`ファイルを編集
```bash
cd /Users/tomo/Local/other/gRPC/service/auth/docker/
vi .env
```

以下のように編集：
```env
# Gmail App Password（生成した16桁のパスワードを入力）
GMAIL_APP_PASSWORD=abcd-efgh-ijkl-mnop  # スペースをハイフンに変更
```

## 3. Docker Composeの起動

```bash
cd /Users/tomo/Local/other/gRPC/service/auth/docker/
docker-compose up -d
```

## 4. テスト

ユーザー登録のテストを実行して、実際にメールが送信されることを確認：

```bash
cd /Users/tomo/Local/other/gRPC/service/auth/app/
go run cmd/client/main.go
```

## トラブルシューティング 🔧

### よくあるエラー

1. **"Username and Password not accepted"**
   - アプリパスワードが正しくない
   - 2段階認証が有効になっていない

2. **"Connection refused"**
   - インターネット接続を確認
   - ファイアウォール設定を確認

3. **"Message rejected"**
   - 送信元アドレスがGmailアカウントと一致していない

### デバッグ方法

ログを確認：
```bash
docker-compose logs auth_service
```

## セキュリティ注意事項 ⚠️

- `.env`ファイルは絶対にGitにコミットしない
- 本番環境では更に強固な認証方式を検討する
- アプリパスワードは定期的に更新する
