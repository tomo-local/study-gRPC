package service

import (
	"context"
	"fmt"
	"time"

	"auth/auth"
	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ResetPassword はパスワードのリセットを行います
func (s *authServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// パスワードの形式チェック
	if !auth.IsValidPassword(req.NewPassword) {
		return &pb.ResetPasswordResponse{
			Success: false,
			Message: "password must be at least 8 characters long",
		}, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	err := s.db.StartTransaction(func(tx *gorm.DB) error {
		// トークンを取得
		resetToken, err := s.db.GetPasswordResetTokenByToken(tx, req.Token)
		if err != nil {
			return status.Error(codes.Unauthenticated, "invalid reset token")
		}

		// トークンが期限切れでないかチェック
		if resetToken.IsExpired() {
			return status.Error(codes.Unauthenticated, "reset token has expired")
		}

		// トークンが使用済みでないかチェック
		if resetToken.Used {
			return status.Error(codes.Unauthenticated, "reset token has already been used")
		}

		// ユーザーを取得
		user, err := s.db.GetUserByID(tx, resetToken.UserID)
		if err != nil {
			return status.Error(codes.Internal, "user not found")
		}

		// 新しいパスワードをハッシュ化
		hashedPassword, err := auth.HashPassword(req.NewPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// ユーザーのパスワードを更新
		user.PasswordHash = hashedPassword
		user.UpdatedAt = time.Now()

		if err := s.db.UpdateUser(tx, user); err != nil {
			return fmt.Errorf("failed to update user password: %w", err)
		}

		// トークンを使用済みにマーク
		resetToken.Used = true
		if err := s.db.UpdatePasswordResetToken(tx, resetToken); err != nil {
			return fmt.Errorf("failed to update reset token: %w", err)
		}

		return nil
	})

	if err != nil {
		return &pb.ResetPasswordResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.ResetPasswordResponse{
		Success: true,
		Message: "Password reset successfully.",
	}, nil
}
