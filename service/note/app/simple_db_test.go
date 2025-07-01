package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("🧪 PostgreSQL Simple Connection Test")

	// データベース接続文字列
	dsn := "host=localhost port=5432 user=noteuser password=notepass dbname=notedb sslmode=disable"

	// データベースに接続
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ データベース接続エラー: %v", err)
	}
	defer db.Close()

	// 接続をテスト
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ データベースpingエラー: %v", err)
	}

	log.Println("✅ PostgreSQL接続成功!")

	// notesテーブルの確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&count)
	if err != nil {
		log.Printf("⚠️  notesテーブル読み取りエラー: %v", err)
	} else {
		log.Printf("📋 notesテーブルのレコード数: %d", count)
	}

	// サンプルレコードを表示
	rows, err := db.Query("SELECT id, title, content FROM notes LIMIT 3")
	if err != nil {
		log.Printf("⚠️  レコード取得エラー: %v", err)
		return
	}
	defer rows.Close()

	log.Println("📝 サンプルレコード:")
	for rows.Next() {
		var id, title, content string
		if err := rows.Scan(&id, &title, &content); err != nil {
			log.Printf("⚠️  レコードスキャンエラー: %v", err)
			continue
		}
		fmt.Printf("  ID: %s, Title: %s, Content: %s\n", id, title, content[:min(50, len(content))])
	}

	log.Println("🎉 データベーステスト完了!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
