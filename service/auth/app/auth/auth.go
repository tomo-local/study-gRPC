package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

// GenerateToken はJWTトークンを生成します
func (manager *JWTManager) GenerateToken(userID, email, name string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(manager.secretKey)
}

// VerifyToken はJWTトークンを検証します
func (manager *JWTManager) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}
		return manager.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// HashPassword はパスワードをハッシュ化します
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword はパスワードを検証します
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateRandomToken はランダムなトークンを生成します
func GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValidEmail はメールアドレスの形式をチェックします
func IsValidEmail(email string) bool {
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

// IsValidPassword はパスワードの形式をチェックします
func IsValidPassword(password string) bool {
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
