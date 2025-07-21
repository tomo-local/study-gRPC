package main

import (
	"context"
	"log"
	"os"

	auth_pb "auth/grpc/api"
	note_pb "note/grpc" // Note serviceのprotoが必要

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// getEnv gets environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 環境変数から接続先を取得（デフォルト値付き）
	authServiceURL := getEnv("AUTH_SERVICE_URL", "localhost:9001")
	noteServiceURL := getEnv("NOTE_SERVICE_URL", "localhost:9002")

	log.Printf("🔗 Auth Service: %s", authServiceURL)
	log.Printf("📝 Note Service: %s", noteServiceURL)

	// 1. Authサービスに接続してログイン
	authConn, err := grpc.NewClient(authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	authClient := auth_pb.NewAuthServiceClient(authConn)

	// ログイン
	loginResp, err := authClient.Login(context.Background(), &auth_pb.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Printf("Login successful! Token: %s", loginResp.AccessToken)

	// 2. Noteサービスに接続
	noteConn, err := grpc.NewClient(noteServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to note service: %v", err)
	}
	defer noteConn.Close()

	noteClient := note_pb.NewNoteServiceClient(noteConn)

	// JWTトークンをメタデータとして送信
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+loginResp.AccessToken)

	// ノートを作成
	noteResp, err := noteClient.CreateNote(ctx, &note_pb.CreateNoteRequest{
		Title:    "Auth後の最初のノート",
		Content:  "認証後にノートを作成できました！",
		Category: "diary",
		Tags:     []string{"auth", "success"},
	})
	if err != nil {
		log.Fatalf("Note creation failed: %v", err)
	}

	log.Printf("Note created: %+v", noteResp.Note)
}
