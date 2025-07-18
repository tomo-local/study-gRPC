package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	config "note/config/server"
	"note/db"
	pb "note/grpc"
	"note/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("🚀 Note gRPC Server Starting...")

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	log.Printf("✅ Config loaded - Port: %d, Env: %s", cfg.Port, cfg.Env)

	database, err := db.New(db.Config{
		Host:     cfg.Database.Host,
		Name:     cfg.Database.DBName,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Port:     cfg.Database.Port,
		SSLMode:  cfg.Database.SSLMode,
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully!")

	// gRPCサーバーの設定
	port := cfg.Port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("❌ Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Note サービスを登録
	noteServer := service.NewNoteServer(database)
	pb.RegisterNoteServiceServer(grpcServer, noteServer)

	// リフレクションを有効化（開発環境用）
	if cfg.Env == "development" {
		reflection.Register(grpcServer)
		log.Println("✅ gRPC reflection enabled for development")
	}

	// Graceful shutdown の設定
	go func() {
		log.Printf("🎯 Note gRPC Server listening on port %d", port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ Failed to serve gRPC server: %v", err)
		}
	}()

	// シグナルハンドリング
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Note gRPC Server...")
	grpcServer.GracefulStop()
	log.Println("✅ Note gRPC Server stopped gracefully")
}
