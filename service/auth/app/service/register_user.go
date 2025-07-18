package service

import (
	"context"
	"fmt"
	"time"

	"auth/auth"
	"auth/db/model"
	pb "auth/grpc/api"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Register はユーザー登録を行います
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 入力バリデーション
	if !auth.IsValidEmail(req.Email) {
		return &pb.RegisterResponse{
			Success: false,
			Message: "invalid email format",
		}, status.Error(codes.InvalidArgument, "invalid email format")
	}

	if !auth.IsValidPassword(req.Password) {
		return &pb.RegisterResponse{
			Success: false,
			Message: "password must be at least 8 characters long",
		}, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	if req.Name == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "name is required",
		}, status.Error(codes.InvalidArgument, "name is required")
	}

	// パスワードをハッシュ化
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "failed to hash password",
		}, status.Error(codes.Internal, "failed to hash password")
	}

	var user *model.User
	err = s.db.StartTransaction(func(tx *gorm.DB) error {
		// ユーザーが既に存在するかチェック
		_, err := s.db.GetUserByEmail(tx, req.Email)
		if err == nil {
			return fmt.Errorf("user already exists with email: %s", req.Email)
		}

		// 新しいユーザーを作成
		user = &model.User{
			ID:            uuid.New().String(),
			Email:         req.Email,
			Name:          req.Name,
			PasswordHash:  hashedPassword,
			EmailVerified: false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.db.CreateUser(tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// メール認証トークンを生成
		token, err := auth.GenerateRandomToken()
		if err != nil {
			return fmt.Errorf("failed to generate verification token: %w", err)
		}

		// メール認証トークンを保存
		emailToken := &model.EmailVerificationToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour), // 24時間有効
			CreatedAt: time.Now(),
		}

		if err := s.db.CreateEmailVerificationToken(tx, emailToken); err != nil {
			return fmt.Errorf("failed to create email verification token: %w", err)
		}

		// 認証メールを送信
		if err := s.mailer.SendVerificationEmail(req.Email, req.Name, token); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}

		return nil
	})

	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully. Please check your email to verify your account.",
		UserId:  user.ID,
	}, nil
}
