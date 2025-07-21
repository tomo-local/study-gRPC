package service

import (
	"context"
	"fmt"
	"time"

	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ResetPassword はパスワードのリセットを行います
func (s *authServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	// 入力の検証
	if err := s.validateResetPasswordInput(req); err != nil {
		return nil, err
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
		hashedPassword, err := hashPassword(req.NewPassword)
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ResetPasswordResponse{}, nil
}

func (s *authServer) validateResetPasswordInput(req *pb.ResetPasswordRequest) error {
	if !isValidPassword(req.NewPassword) {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	return nil
}
