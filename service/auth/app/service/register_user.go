package service

import (
	"context"
	"fmt"
	"time"

	"auth/db/model"
	pb "auth/grpc/api"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Register はユーザー登録を行います
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 入力の検証
	if err := s.validateRegistrationInput(req); err != nil {
		return nil, err
	}

	// パスワードをハッシュ化
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	var user *model.User
	err = s.db.StartTransaction(func(tx *gorm.DB) error {
		// 新しいユーザーを作成
		user = &model.User{
			ID:            uuid.New().String(),
			Email:         req.Email,
			Name:          req.Name,
			PasswordHash:  hashedPassword,
			EmailVerified: false,
		}

		if err := s.db.CreateUser(tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// メール認証トークンを生成
		token, err := generateRandomToken()
		if err != nil {
			return fmt.Errorf("failed to generate verification token: %w", err)
		}

		// メール認証トークンを保存
		emailToken := &model.EmailVerificationToken{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour), // 24時間有効
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		UserId: user.ID,
	}, nil
}

func (s *authServer) validateRegistrationInput(req *pb.RegisterRequest) error {
	if !isValidEmail(req.Email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	if !isValidPassword(req.Password) {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	if user, err := s.db.GetUserByEmail(nil, req.Email); err != nil {
		// ユーザーが存在しない場合はnilを返す
		return status.Error(codes.Internal, "failed to check existing user")
	} else if user != nil {
		return status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	return nil
}
