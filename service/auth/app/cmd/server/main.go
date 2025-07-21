package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth/auth"
	config "auth/config/server"
	"auth/db"
	pb "auth/grpc/api"
	"auth/mailer"
	"auth/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 設定を読み込み（パスは自動決定）
	cfg, err := config.LoadConfigWithAutoPath()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続を初期化
	dbConfig := db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}

	database, err := db.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// データベース接続をテスト
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// JWTマネージャーを初期化
	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey, time.Duration(cfg.JWT.TokenDurationH)*time.Hour)

	// メール送信機能を初期化
	mailerConfig := mailer.Config{
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
		From:     cfg.Mailer.From,
	}
	mailerService := mailer.New(mailerConfig)

	// 認証サービスを初期化
	authService := service.NewAuthServer(database, jwtManager, mailerService)

	// gRPCサーバーを作成
	grpcServer := grpc.NewServer()

	// サービスを登録
	pb.RegisterAuthServiceServer(grpcServer, authService)

	// リフレクション機能を有効化（開発時に便利）
	reflection.Register(grpcServer)

	// リスナーを作成
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Auth gRPC server is running on port %d", cfg.Server.Port)

	// Graceful shutdown のためのチャンネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// サーバーをゴルーチンで起動
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// シャットダウンシグナルを待機
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	grpcServer.GracefulStop()
	database.Close()
	log.Println("Server stopped")
}
