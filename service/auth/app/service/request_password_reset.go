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

// RequestPasswordReset はパスワードリセットの要求を処理します
func (s *authServer) RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error) {
	err := s.db.StartTransaction(func(tx *gorm.DB) error {
		// ユーザーを取得
		user, err := s.db.GetUserByEmail(tx, req.Email)
		if err != nil {
			// セキュリティ上の理由で、ユーザーが存在しない場合でも成功として扱う
			return nil
		}

		// パスワードリセットトークンを生成
		token, err := auth.GenerateRandomToken()
		if err != nil {
			return fmt.Errorf("failed to generate reset token: %w", err)
		}

		// パスワードリセットトークンを保存
		resetToken := &model.PasswordResetToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(1 * time.Hour), // 1時間有効
			CreatedAt: time.Now(),
		}

		if err := s.db.CreatePasswordResetToken(tx, resetToken); err != nil {
			return fmt.Errorf("failed to create password reset token: %w", err)
		}

		// パスワードリセットメールを送信
		if err := s.mailer.SendPasswordResetEmail(req.Email, user.Name, token); err != nil {
			return fmt.Errorf("failed to send password reset email: %w", err)
		}

		return nil
	})

	if err != nil {
		return &pb.RequestPasswordResetResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.RequestPasswordResetResponse{
		Success: true,
		Message: "Password reset email sent successfully.",
	}, nil
}
