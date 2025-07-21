package main

import (
	"fmt"
	"log"
	"os"

	config "auth/config/server"
)

func main() {
	fmt.Println("=== Config Test ===")

	// 環境変数のテスト
	os.Setenv("DB_HOST", "test-host-from-env")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("JWT_SECRET_KEY", "test-secret-from-env")

	fmt.Println("Environment variables set:")
	fmt.Printf("  DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("  SERVER_PORT: %s\n", os.Getenv("SERVER_PORT"))
	fmt.Printf("  JWT_SECRET_KEY: %s\n", os.Getenv("JWT_SECRET_KEY"))

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("\n=== Loaded Config ===")
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("Database Port: %d\n", cfg.Database.Port)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("Database Name: %s\n", cfg.Database.Name)
	fmt.Printf("JWT Secret: %s\n", cfg.JWT.SecretKey)
	fmt.Printf("JWT Duration: %d hours\n", cfg.JWT.TokenDurationH)
	fmt.Printf("Mailer Host: %s\n", cfg.Mailer.Host)

	fmt.Println("\n✅ Config test completed successfully!")
}
