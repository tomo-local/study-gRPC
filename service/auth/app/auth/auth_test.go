package auth

import (
	"testing"
	"time"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"", false},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"test@.com", false},
		{"test@com", false},
		{"test..test@example.com", false},
		{".test@example.com", false},
		{"test@example.com.", false},
		{"test@example..com", false},
		{"test@example@com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"password123", true},
		{"Password1", true},
		{"mypassword123", true},
		{"", false},
		{"short", false},
		{"12345678", false},                // 数字のみ
		{"password", false},                // 文字のみ
		{"PASSWORD123", true},              // 大文字含む
		{"password123!", true},             // 特殊文字含む
		{"a1", false},                      // 短すぎる
		{string(make([]byte, 129)), false}, // 長すぎる
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := IsValidPassword(tt.password)
			if result != tt.expected {
				t.Errorf("IsValidPassword(%q) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty string")
	}

	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 正しいパスワードのテスト
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword failed for correct password")
	}

	// 間違ったパスワードのテスト
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword succeeded for wrong password")
	}
}

func TestJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Hour

	manager := NewJWTManager(secretKey, tokenDuration)

	userID := "test-user-id"
	email := "test@example.com"
	name := "Test User"

	// トークン生成のテスト
	token, err := manager.GenerateToken(userID, email, name)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}

	// トークン検証のテスト
	claims, err := manager.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}

	if claims.Name != name {
		t.Errorf("Expected Name %s, got %s", name, claims.Name)
	}

	// 無効なトークンのテスト
	_, err = manager.VerifyToken("invalid-token")
	if err == nil {
		t.Error("VerifyToken should fail for invalid token")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token1, err := GenerateRandomToken()
	if err != nil {
		t.Fatalf("GenerateRandomToken failed: %v", err)
	}

	token2, err := GenerateRandomToken()
	if err != nil {
		t.Fatalf("GenerateRandomToken failed: %v", err)
	}

	if token1 == token2 {
		t.Error("GenerateRandomToken should generate different tokens")
	}

	if len(token1) == 0 {
		t.Error("GenerateRandomToken returned empty token")
	}

	// トークンの長さは64文字（32バイト * 2）
	if len(token1) != 64 {
		t.Errorf("Expected token length 64, got %d", len(token1))
	}
}
