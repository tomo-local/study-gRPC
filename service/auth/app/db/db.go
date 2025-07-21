package db

import (
	"fmt"
	"time"

	"auth/db/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	StartTransaction(fn func(tx *gorm.DB) error) error
	Close() error
	Ping() error

	// User関連のメソッド
	CreateUser(tx *gorm.DB, user *model.User) error
	GetUserByID(tx *gorm.DB, id string) (*model.User, error)
	GetUserByEmail(tx *gorm.DB, email string) (*model.User, error)
	UpdateUser(tx *gorm.DB, user *model.User) error
	DeleteUser(tx *gorm.DB, id string) error

	// EmailVerificationToken関連のメソッド
	CreateEmailVerificationToken(tx *gorm.DB, token *model.EmailVerificationToken) error
	GetEmailVerificationTokenByToken(tx *gorm.DB, token string) (*model.EmailVerificationToken, error)
	UpdateEmailVerificationToken(tx *gorm.DB, token *model.EmailVerificationToken) error
	DeleteEmailVerificationToken(tx *gorm.DB, id string) error

	// PasswordResetToken関連のメソッド
	CreatePasswordResetToken(tx *gorm.DB, token *model.PasswordResetToken) error
	GetPasswordResetTokenByToken(tx *gorm.DB, token string) (*model.PasswordResetToken, error)
	UpdatePasswordResetToken(tx *gorm.DB, token *model.PasswordResetToken) error
	DeletePasswordResetToken(tx *gorm.DB, id string) error
	DeleteExpiredPasswordResetTokens(tx *gorm.DB) error
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

// migrateDB はデータベースのマイグレーションを実行します。
func migrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.EmailVerificationToken{},
		&model.PasswordResetToken{},
	); err != nil {
		return err
	}

	return nil
}

// Close は終了時にデータベース接続を閉じます。
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

// Ping はデータベースへの接続が生きているかを確認します。
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
