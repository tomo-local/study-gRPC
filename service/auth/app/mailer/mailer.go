package mailer

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func New(config Config) *Mailer {
	return &Mailer{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
	}
}

// SendVerificationEmail はメールアドレス認証用のメールを送信します
func (m *Mailer) SendVerificationEmail(to, name, token string) error {
	subject := "メールアドレス認証のお願い"
	body := fmt.Sprintf(`
%s 様

アカウント登録ありがとうございます。

以下のリンクをクリックしてメールアドレスの認証を完了してください。

認証URL: http://localhost:3000/verify-email?token=%s

このリンクは24時間有効です。

よろしくお願いいたします。
`, name, token)

	return m.sendEmail(to, subject, body)
}

// SendPasswordResetEmail はパスワードリセット用のメールを送信します
func (m *Mailer) SendPasswordResetEmail(to, name, token string) error {
	subject := "パスワードリセットのお知らせ"
	body := fmt.Sprintf(`
%s 様

パスワードリセットのリクエストを受け付けました。

以下のリンクをクリックしてパスワードの再設定を行ってください。

パスワードリセットURL: http://localhost:3000/reset-password?token=%s

このリンクは1時間有効です。

身に覚えのないリクエストの場合は、このメールを無視してください。

よろしくお願いいたします。
`, name, token)

	return m.sendEmail(to, subject, body)
}

// sendEmail はメールを送信します
func (m *Mailer) sendEmail(to, subject, body string) error {
	// メール設定が空の場合はログ出力のみ（開発環境用）
	if m.username == "" || m.password == "" {
		log.Printf("Email would be sent to %s with subject: %s", to, subject)
		log.Printf("Email body: %s", body)
		return nil
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	d := gomail.NewDialer(m.host, m.port, m.username, m.password)

	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// IsValidEmailFormat はメールアドレスの形式をチェックします
func IsValidEmailFormat(email string) bool {
	// 基本的なメールアドレスの形式チェック
	return len(email) > 0 &&
		len(email) <= 320 &&
		strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		!strings.HasPrefix(email, "@") &&
		!strings.HasSuffix(email, "@") &&
		!strings.Contains(email, "..") &&
		strings.Count(email, "@") == 1
}
