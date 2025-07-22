package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// Auth Serviceの proto定義（実際のパスに合わせて調整）
	// auth_pb "auth/grpc/api"
)

// AuthClient はAuth Serviceとの通信を管理するクライアント
type AuthClient struct {
	conn *grpc.ClientConn
	// client     auth_pb.AuthServiceClient
	serviceURL string
	timeout    time.Duration
	enabled    bool
}

// NewAuthClient creates a new Auth Service client
func NewAuthClient(serviceURL string, timeout time.Duration, enabled bool) (*AuthClient, error) {
	if !enabled {
		log.Println("🔒 Auth Service integration is disabled")
		return &AuthClient{
			serviceURL: serviceURL,
			timeout:    timeout,
			enabled:    false,
		}, nil
	}

	log.Printf("🔐 Connecting to Auth Service: %s", serviceURL)

	conn, err := grpc.NewClient(serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", serviceURL, err)
	}

	// client := auth_pb.NewAuthServiceClient(conn)

	log.Printf("✅ Successfully connected to Auth Service")

	return &AuthClient{
		conn: conn,
		// client:     client,
		serviceURL: serviceURL,
		timeout:    timeout,
		enabled:    true,
	}, nil
}

// Close closes the connection to Auth Service
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsEnabled returns whether auth service integration is enabled
func (c *AuthClient) IsEnabled() bool {
	return c.enabled
}

// VerifyToken verifies JWT token with Auth Service
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*UserInfo, error) {
	if !c.enabled {
		log.Println("⚠️  Auth verification skipped (disabled)")
		return &UserInfo{
			UserID: "anonymous",
			Email:  "anonymous@example.com",
			Name:   "Anonymous User",
		}, nil
	}

	// 実際の実装では Auth Service の VerifyToken を呼び出し
	/*
		req := &auth_pb.VerifyTokenRequest{
			Token: token,
		}

		resp, err := c.client.VerifyToken(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("token verification failed: %w", err)
		}

		if !resp.Valid {
			return nil, fmt.Errorf("invalid token")
		}

		return &UserInfo{
			UserID: resp.UserId,
			Email:  resp.Email,
			Name:   resp.Name,
		}, nil
	*/

	// 仮実装：実際の開発では上記のコメントアウト部分を使用
	log.Printf("🔍 Verifying token with Auth Service (mock)")
	return &UserInfo{
		UserID: "user123",
		Email:  "user@example.com",
		Name:   "Test User",
	}, nil
}

// UserInfo represents user information from auth service
type UserInfo struct {
	UserID string
	Email  string
	Name   string
}

// HealthCheck checks if Auth Service is available
func (c *AuthClient) HealthCheck(ctx context.Context) error {
	if !c.enabled {
		return nil
	}

	// 簡単な接続確認（実際の実装ではHealth check RPCを使用）
	log.Printf("🏥 Auth Service health check: %s", c.serviceURL)

	// 実装例：
	// req := &health_pb.HealthCheckRequest{Service: "auth"}
	// _, err := c.healthClient.Check(ctx, req)
	// return err

	return nil
}
