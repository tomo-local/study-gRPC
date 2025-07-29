# メールサーバー構築 必須知識ガイド

> **関連文書**:
> - [ADR-001: 学習用メールサーバーのアーキテクチャ選択](./adr-001-mail-server-architecture.md)
> - [Go言語メールサーバー学習計画プラン](./mail-server-learning-plan.md)

## 概要

メールサーバーの構築には、ネットワーク、プロトコル、セキュリティ、システム管理など多岐にわたる知識が必要です。本ガイドでは、学習レベル別に必要な知識を整理し、効率的な学習パスを提示します。

---

## 1. 基礎レベル（必須知識）

### 1.1 ネットワーク基礎

#### **TCP/IP基礎**
```bash
# 理解すべき概念
- IPアドレス（IPv4/IPv6）
- ポート番号の概念
- TCP接続の仕組み
- ファイアウォールの基本
```

**実践例:**
```bash
# SMTPで使用される主要ポート
25   # SMTP（サーバー間通信）
587  # Submission（クライアント送信）
465  # SMTPS（SSL/TLS）
993  # IMAPS
995  # POP3S
```

#### **DNS（ドメインネームシステム）**
```dns
# 必須DNSレコードタイプ
A      example.com.           192.168.1.100
MX     example.com.      10   mail.example.com.
TXT    example.com.           "v=spf1 ip4:192.168.1.100 ~all"
CNAME  mail.example.com.      server.example.com.
```

**学習ポイント:**
- MXレコードの優先度
- DNSの伝播時間（TTL）
- nslookup/digコマンドの使用法

### 1.2 メールプロトコル基礎

#### **SMTP（Simple Mail Transfer Protocol）**
```smtp
# 基本的なSMTPコマンド
HELO/EHLO    # サーバー識別
MAIL FROM    # 送信者指定
RCPT TO      # 受信者指定
DATA         # メール本文開始
QUIT         # 接続終了
```

**実践例:**
```bash
# telnetでSMTPテスト
telnet mail.example.com 25
> EHLO client.example.com
> MAIL FROM:<sender@example.com>
> RCPT TO:<recipient@example.com>
> DATA
> Subject: Test
>
> Hello World
> .
> QUIT
```

#### **メールヘッダー構造**
```email
# 重要なヘッダー項目
From: sender@example.com
To: recipient@example.com
Subject: メールサブジェクト
Date: Mon, 1 Jan 2024 12:00:00 +0900
Message-ID: <unique-id@example.com>
Received: by mail.example.com (Postfix)
```

### 1.3 プログラミング基礎（Go言語）

#### **Go基本文法**
```go
// 必須項目
- 変数、定数、データ型
- 関数とメソッド
- 構造体とインターフェース
- エラーハンドリング
- パッケージ管理（go modules）
```

#### **並行処理（Goroutine/Channel）**
```go
// SMTPサーバーでの並行処理例
func handleConnection(conn net.Conn) {
    defer conn.Close()
    // 各接続を個別のgoroutineで処理
}

func main() {
    listener, _ := net.Listen("tcp", ":25")
    for {
        conn, _ := listener.Accept()
        go handleConnection(conn) // 並行処理
    }
}
```

### 1.4 データベース基礎（PostgreSQL）

#### **SQL基本操作**
```sql
-- ユーザー認証クエリ例
SELECT id, password_hash
FROM virtual_users
WHERE email = $1 AND enabled = true;

-- メール保存
INSERT INTO messages (user_id, subject, sender, body)
VALUES ($1, $2, $3, $4);
```

#### **接続管理**
```go
// Go でのPostgreSQL接続
import (
    "database/sql"
    _ "github.com/lib/pq"
)

db, err := sql.Open("postgres",
    "user=mail dbname=mailserver sslmode=disable")
```

---

## 2. 中級レベル（実装に必要）

### 2.1 メールサーバーアーキテクチャ

#### **MTA（Mail Transfer Agent）設計**
```go
// SMTPサーバーの基本構造
type SMTPServer struct {
    hostname string
    port     int
    db       *sql.DB
    auth     AuthHandler
    delivery DeliveryHandler
}

// 認証処理
type AuthHandler interface {
    Authenticate(username, password string) error
}

// 配信処理
type DeliveryHandler interface {
    DeliverLocal(message *Message) error
    DeliverRemote(message *Message) error
}
```

#### **メッセージキューイング**
```go
// 配信キュー管理
type MessageQueue struct {
    pending   chan *Message
    retry     chan *Message
    failed    chan *Message
    workers   int
}

func (mq *MessageQueue) ProcessQueue() {
    for i := 0; i < mq.workers; i++ {
        go mq.worker()
    }
}
```

### 2.2 認証とセキュリティ

#### **SMTP認証メカニズム**
```go
// PLAIN認証の実装例
func (a *AuthHandler) AuthPlain(identity, username, password string) error {
    // Base64デコード
    decoded, _ := base64.StdEncoding.DecodeString(password)

    // データベースでの認証確認
    return a.validateCredentials(username, string(decoded))
}
```

#### **TLS/SSL実装**
```go
// TLS設定
tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    ServerName:   "mail.example.com",
}

// STARTTLS対応
listener := tls.NewListener(netListener, tlsConfig)
```

### 2.3 Docker化

#### **マルチステージビルド**
```dockerfile
# ビルドステージ
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o smtp-server

# 実行ステージ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/smtp-server /smtp-server
EXPOSE 25 587
CMD ["/smtp-server"]
```

#### **Docker Compose構成**
```yaml
version: '3.8'
services:
  smtp-server:
    build: ./smtp-server
    ports:
      - "25:25"
      - "587:587"
    environment:
      - DB_HOST=postgresql
      - SMTP_HOSTNAME=mail.example.com
    depends_on:
      - postgresql

  postgresql:
    image: postgres:13
    environment:
      POSTGRES_DB: mailserver
      POSTGRES_USER: mail
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

---

## 3. 上級レベル（本格運用）

### 3.1 DNS とメール認証

#### **SPF（Sender Policy Framework）**
```dns
# SPF レコード設定
example.com. TXT "v=spf1 ip4:192.168.1.100 include:_spf.google.com ~all"

# SPF 検証機能
- ip4/ip6: 許可IPアドレス
- include: 他ドメインのSPF参照
- ~all: ソフトフェイル（推奨）
- -all: ハードフェイル
```

#### **DKIM（DomainKeys Identified Mail）**
```go
// DKIM署名実装
import "github.com/emersion/go-msgauth/dkim"

// 秘密鍵でメッセージ署名
options := &dkim.SignOptions{
    Domain:   "example.com",
    Selector: "default",
    Signer:   privateKey,
}

signedMessage, err := dkim.Sign(message, options)
```

```dns
# DKIM公開鍵レコード
default._domainkey.example.com. TXT "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA..."
```

#### **DMARC（Domain-based Message Authentication）**
```dns
# DMARC ポリシー
_dmarc.example.com. TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"

# パラメータ説明
- p=none|quarantine|reject  # ポリシー
- rua=                      # 集計レポート送信先
- ruf=                      # フォレンジックレポート送信先
```

### 3.2 パフォーマンスと監視

#### **接続プール管理**
```go
// データベース接続プール
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// SMTP接続制限
type ConnectionLimiter struct {
    semaphore chan struct{}
}

func NewConnectionLimiter(maxConns int) *ConnectionLimiter {
    return &ConnectionLimiter{
        semaphore: make(chan struct{}, maxConns),
    }
}
```

#### **ログ記録と監視**
```go
// 構造化ログ
import "github.com/sirupsen/logrus"

log.WithFields(logrus.Fields{
    "client_ip": clientIP,
    "from":      fromAddress,
    "to":        toAddress,
    "size":      messageSize,
}).Info("Message accepted")

// メトリクス収集
type Metrics struct {
    messagesReceived counter
    messagesDelivered counter
    authFailures     counter
    connectionCount  gauge
}
```

### 3.3 高可用性とスケーリング

#### **負荷分散**
```yaml
# HAProxy設定例
backend smtp_servers
    balance roundrobin
    server smtp1 192.168.1.10:25 check
    server smtp2 192.168.1.11:25 check
```

#### **データベース冗長化**
```yaml
# PostgreSQL Master-Slave構成
services:
  postgres-master:
    image: postgres:13
    environment:
      POSTGRES_REPLICATION_USER: replicator

  postgres-slave:
    image: postgres:13
    environment:
      PGUSER: replicator
      PGPASSWORD: password
      POSTGRES_MASTER_SERVICE: postgres-master
```

---

## 4. 学習リソースと参考資料

### 4.1 RFC文書（必読）
- **[RFC 5321](https://tools.ietf.org/html/rfc5321)** - SMTP Protocol
- **[RFC 5322](https://tools.ietf.org/html/rfc5322)** - Internet Message Format
- **[RFC 7208](https://tools.ietf.org/html/rfc7208)** - SPF
- **[RFC 6376](https://tools.ietf.org/html/rfc6376)** - DKIM
- **[RFC 7489](https://tools.ietf.org/html/rfc7489)** - DMARC

### 4.2 Go言語ライブラリ
```go
// 主要ライブラリ
"github.com/emersion/go-smtp"      // SMTP実装
"github.com/emersion/go-imap"      // IMAP実装
"github.com/emersion/go-msgauth"   // SPF/DKIM/DMARC
"github.com/jordan-wright/email"   // メール作成支援
"github.com/lib/pq"                // PostgreSQL
"gorm.io/gorm"                     // ORM
```

### 4.3 テスト・デバッグツール
```bash
# コマンドラインツール
telnet          # SMTP手動テスト
openssl s_client # TLS接続テスト
dig/nslookup    # DNS確認
swaks           # SMTP テストツール

# オンラインツール
- MXToolbox (DNS/SPF/DKIM チェック)
- Mail-tester (メール品質スコア)
- DMARC Analyzer
```

---

## 5. 学習パス提案

### 🏁 **Phase 0: 事前学習（1週間）**
1. TCP/IP, DNS基礎の復習
2. Go言語基本文法の確認
3. Docker基本操作の習得
4. PostgreSQL基本操作の確認

### 📚 **Phase 1: 基礎実装（2週間）**
1. SMTP プロトコルの理解
2. 基本的なGo SMTPサーバー実装
3. PostgreSQL認証連携
4. Docker環境構築

### 🔧 **Phase 2: 実用実装（3週間）**
1. 実ドメイン・DNS設定
2. TLS/SSL実装
3. SPF/DKIM実装
4. 監視・ログ機能

### 🚀 **Phase 3: 本格運用（2週間）**
1. Gmail送信実現
2. DMARC実装
3. パフォーマンス最適化
4. 運用手順の確立

---

## 6. よくある躓きポイントと対策

### ❌ **技術的な躓きポイント**

| 問題 | 原因 | 対策 |
|------|------|------|
| SMTP認証に失敗 | Base64エンコード誤解 | RFC仕様の正確な理解 |
| DNS設定が反映されない | TTL設定、キャッシュ | 伝播時間待機、確認コマンド |
| TLS接続エラー | 証明書設定誤り | Let's Encrypt で検証 |
| Gmail送信でスパム判定 | SPF/DKIM未設定 | 認証設定の完全実装 |

### 💡 **学習上の躓きポイント**

| 問題 | 対策 |
|------|------|
| プロトコル仕様が複雑 | 段階的実装、動作確認 |
| Go言語のクセ | サンプルコード活用 |
| Docker設定の複雑さ | 最小構成から開始 |
| 全体像の把握困難 | 図解、アーキテクチャ整理 |

---

## 7. 成功のためのポイント

### ✅ **推奨学習方法**
1. **理論と実践のバランス** - RFC読解と実装を並行
2. **段階的な複雑化** - 動作するものから機能追加
3. **ログとテストの重視** - 問題の早期発見
4. **コミュニティ活用** - 質問、情報収集

### ⚡ **効率化のコツ**
- **既存ライブラリの活用** - ゼロから実装しない
- **Docker活用** - 環境依存の問題を回避
- **自動テストの実装** - 回帰テスト確保
- **段階的デプロイ** - リスク最小化

---

**学習開始前チェックリスト:**
- [ ] Go言語の基本文法は理解している
- [ ] Docker の基本操作ができる
- [ ] PostgreSQL の基本操作ができる
- [ ] TCP/IP, DNS の基礎知識がある
- [ ] 学習に充てる時間を確保している

> 💡 **重要**: 完璧な知識を身につけてから始める必要はありません。実装しながら学習することで、より深い理解が得られます。
