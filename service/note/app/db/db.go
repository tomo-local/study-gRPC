package db

import (
	"fmt"
	"note/db/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	StartTransaction(fn func(tx *gorm.DB) error) error
	Close() error
	Ping() error
	// Note関連のメソッド

	CreateNote(tx *gorm.DB, note *model.Note) error
	UpdateNote(tx *gorm.DB, note *model.Note) error
	DeleteNote(tx *gorm.DB, id string) error
	GetNoteByID(tx *gorm.DB, id string) (*model.Note, error)
	ListNotes(tx *gorm.DB, page, limit int32, category string, tags []string) ([]*model.Note, int64, error)
}

type db struct {
	client *gorm.DB
}

type Config struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     int
	SSLMode  string
}

const retrySleepPeriod = 5 * time.Second

func InitDB(config Config) (*gorm.DB, error) {
	// データベース接続を初期化
	db, err := connectDB(config, 5)
	if err != nil {
		return nil, err
	}

	// マイグレーションを実行
	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB(cfg Config, retryCount int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Tokyo",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if retryCount > 0 {
			time.Sleep(retrySleepPeriod)
			return connectDB(cfg, retryCount-1)
		}

		return nil, err
	}

	return db, nil
}

func New(cfg Config) (DB, error) {
	client, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &db{
		client: client,
	}, nil
}

// Migrate はデータベースのマイグレーションを実行します。
func migrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.Note{},
	); err != nil {
		return err
	}

	return nil
}

// 終了時にデータベース接続を閉じます。
// これは、アプリケーションのシャットダウン時に呼び出すことを想定しています。
func (d *db) Close() error {
	sqlDB, err := d.client.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// データベースへの接続が生きているかを確認します。
func (d *db) Ping() error {
	sqlDB, err := d.client.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
