package main

import (
	"context"
	"log"
	"time"

	auth_pb "auth/grpc/api" // Authサービス
	note_pb "note/grpc"     // Noteサービス

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("🚀 マイクロサービス連携デモ開始")

	// Step 1: Authサービスでユーザー登録とログイン
	accessToken, err := authenticateUser()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Step 2: 認証トークンを使ってNoteサービスでノートを作成
	err = createNoteWithAuth(accessToken)
	if err != nil {
		log.Fatalf("Note creation failed: %v", err)
	}

	log.Println("✅ マイクロサービス連携完了！")
}

// authenticateUser はAuthサービスでユーザーをログインさせる
func authenticateUser() (string, error) {
	log.Println("📝 Authサービスに接続中...")

	// Docker Composeではサービス名で接続可能
	// ローカル開発時は localhost:9001 (外側ポート)
	authConn, err := grpc.NewClient("localhost:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer authConn.Close()

	authClient := auth_pb.NewAuthServiceClient(authConn)
	ctx := context.Background()

	// まずユーザー登録を試行
	log.Println("👤 ユーザー登録中...")
	registerReq := &auth_pb.RegisterRequest{
		Email:    "demo@example.com",
		Password: "password123",
		Name:     "Demo User",
	}

	registerResp, err := authClient.Register(ctx, registerReq)
	if err != nil {
		log.Printf("Registration failed (might already exist): %v", err)
	} else {
		log.Printf("Registration: %+v", registerResp)
	}

	// ログイン
	log.Println("🔐 ログイン中...")
	loginReq := &auth_pb.LoginRequest{
		Email:    "demo@example.com",
		Password: "password123",
	}

	loginResp, err := authClient.Login(ctx, loginReq)
	if err != nil {
		return "", err
	}

	log.Printf("✅ ログイン成功! ユーザー: %s", loginResp.User.Name)
	return loginResp.AccessToken, nil
}

// createNoteWithAuth は認証トークンを使ってNoteサービスでノートを作成
func createNoteWithAuth(accessToken string) error {
	log.Println("📄 Noteサービスに接続中...")

	// Docker Composeではサービス名で接続可能
	// ローカル開発時は localhost:9002 (外側ポート)
	noteConn, err := grpc.NewClient("localhost:9002",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer noteConn.Close()

	noteClient := note_pb.NewNoteServiceClient(noteConn)

	// JWTトークンをAuthorizationヘッダーとして送信
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accessToken)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// ノートを作成
	log.Println("📝 認証後のノート作成中...")
	createReq := &note_pb.CreateNoteRequest{
		Title:    "認証後の最初のノート",
		Content:  "Auth Service経由でログインした後、Note Serviceでノートを作成しました！",
		Category: "microservices",
		Tags:     []string{"auth", "grpc", "microservices", "success"},
	}

	createResp, err := noteClient.CreateNote(ctx, createReq)
	if err != nil {
		return err
	}

	log.Printf("✅ ノート作成成功!")
	log.Printf("   ID: %s", createResp.Note.Id)
	log.Printf("   Title: %s", createResp.Note.Title)
	log.Printf("   Category: %s", createResp.Note.Category)
	log.Printf("   Tags: %v", createResp.Note.Tags)

	// 作成したノートを取得して確認
	getReq := &note_pb.GetNoteRequest{
		Id: createResp.Note.Id,
	}

	getResp, err := noteClient.GetNote(ctx, getReq)
	if err != nil {
		log.Printf("Note retrieval failed: %v", err)
		return err
	}

	log.Println("📖 作成したノートの内容確認:")
	log.Printf("   Content: %s", getResp.Note.Content)
	log.Printf("   Created At: %s", getResp.Note.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))

	return nil
}
