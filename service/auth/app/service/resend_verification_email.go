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

// ResendVerificationEmail は認証メールを再送信します
func (s *authServer) ResendVerificationEmail(ctx context.Context, req *pb.ResendVerificationEmailRequest) (*pb.ResendVerificationEmailResponse, error) {
	err := s.db.StartTransaction(func(tx *gorm.DB) error {
		// ユーザーを取得
		user, err := s.db.GetUserByEmail(tx, req.Email)
		if err != nil {
			return status.Error(codes.NotFound, "user not found")
		}

		// 既に認証済みかチェック
		if user.EmailVerified {
			return status.Error(codes.AlreadyExists, "email already verified")
		}

		// 新しいメール認証トークンを生成
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
		if err := s.mailer.SendVerificationEmail(req.Email, user.Name, token); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}

		return nil
	})

	if err != nil {
		return &pb.ResendVerificationEmailResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.ResendVerificationEmailResponse{
		Success: true,
		Message: "Verification email sent successfully.",
	}, nil
}
