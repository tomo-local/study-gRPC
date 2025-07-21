package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"auth/db/model"
	pb "auth/grpc/api"

	"golang.org/x/crypto/bcrypt"
)

// convertUserToProto はmodel.UserをProtoのUserに変換します
func convertUserToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:            user.ID, // UUIDはstringなのでそのまま使用
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
	}
}

// hashPassword はパスワードをハッシュ化します
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword はパスワードを検証します
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generateRandomToken はランダムなトークンを生成します
func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// isValidEmail はメールアドレスの形式をチェックします
func isValidEmail(email string) bool {
	// より厳密なメールアドレスの検証
	if len(email) == 0 || len(email) > 320 {
		return false
	}

	// @ の個数チェック
	atCount := 0
	for _, c := range email {
		if c == '@' {
			atCount++
		}
	}
	if atCount != 1 {
		return false
	}

	// 基本的な形式チェック
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// ローカル部のチェック
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// ドメイン部のチェック
	if len(domainPart) == 0 || len(domainPart) > 255 {
		return false
	}

	// ドメイン部にドットが含まれているかチェック
	if !strings.Contains(domainPart, ".") {
		return false
	}

	// 先頭・末尾のドットチェック
	if strings.HasPrefix(email, ".") || strings.HasSuffix(email, ".") {
		return false
	}

	// 連続ドットチェック
	if strings.Contains(email, "..") {
		return false
	}

	return true
}

// isValidPassword はパスワードの形式をチェックします
func isValidPassword(password string) bool {
	// 最低8文字、最大128文字
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	// 少なくとも1つの文字と1つの数字を含む
	hasLetter := false
	hasDigit := false

	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
