package main

import (
	"context"
	"log"
	"time"

	pb "auth/grpc/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPCサーバーに接続
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// クライアントを作成
	client := pb.NewAuthServiceClient(conn)

	// コンテキストを作成（タイムアウト付き）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. ユーザー登録のテスト
	log.Println("Testing user registration...")
	registerReq := &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	registerResp, err := client.Register(ctx, registerReq)
	if err != nil {
		log.Printf("Register failed: %v", err)
	} else {
		log.Printf("Register response: %+v", registerResp)
	}

	// 2. ログインのテスト（メール認証前なので失敗するはず）
	log.Println("\nTesting login (should fail - email not verified)...")
	loginReq := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		log.Printf("Login failed (expected): %v", err)
	} else {
		log.Printf("Login response: %+v", loginResp)
	}

	// 3. 認証メール再送のテスト
	log.Println("\nTesting resend verification email...")
	resendReq := &pb.ResendVerificationEmailRequest{
		Email: "test@example.com",
	}

	resendResp, err := client.ResendVerificationEmail(ctx, resendReq)
	if err != nil {
		log.Printf("Resend verification email failed: %v", err)
	} else {
		log.Printf("Resend verification email response: %+v", resendResp)
	}

	// 4. パスワードリセット要求のテスト
	log.Println("\nTesting password reset request...")
	resetReq := &pb.RequestPasswordResetRequest{
		Email: "test@example.com",
	}

	resetResp, err := client.RequestPasswordReset(ctx, resetReq)
	if err != nil {
		log.Printf("Password reset request failed: %v", err)
	} else {
		log.Printf("Password reset request response: %+v", resetResp)
	}

	// 5. 無効なトークンでのトークン検証テスト
	log.Println("\nTesting token verification with invalid token...")
	verifyReq := &pb.VerifyTokenRequest{
		Token: "invalid-token",
	}

	verifyResp, err := client.VerifyToken(ctx, verifyReq)
	if err != nil {
		log.Printf("Token verification failed (expected): %v", err)
	} else {
		log.Printf("Token verification response: %+v", verifyResp)
	}

	log.Println("\nClient test completed!")
}
