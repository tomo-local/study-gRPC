package jwt

import (
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Hour

	manager := NewManager(secretKey, tokenDuration)

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

func TestUserClaims(t *testing.T) {
	userID := "uuid-test-id"
	email := "test@example.com"
	name := "Test User"

	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
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
}
