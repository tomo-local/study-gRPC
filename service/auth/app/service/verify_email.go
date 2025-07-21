package service

import (
	"context"

	"auth/db/model"
	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// VerifyEmail はメールアドレス認証を行います
func (s *authServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// 入力の検証	// トークンを取得
	emailToken, err := s.db.GetEmailVerificationTokenByToken(nil, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get email verification token")
	}

	// トークンが期限切れでないかチェック
	if emailToken.IsExpired() {
		return nil, status.Error(codes.Unauthenticated, "verification token has expired")
	}

	// トークンが使用済みでないかチェック
	if emailToken.Used {
		return nil, status.Error(codes.Unauthenticated, "verification token has already been used")
	}

	// ユーザーを取得
	user, err := s.db.GetUserByID(nil, emailToken.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = s.db.StartTransaction(func(tx *gorm.DB) error {
		updateUser := &model.User{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			PasswordHash:  user.PasswordHash,
			EmailVerified: true,
		}
		if err := s.db.UpdateUser(tx, updateUser); err != nil {
			return status.Error(codes.Internal, "failed to update user")
		}

		updateEmailToken := &model.EmailVerificationToken{
			ID:   emailToken.ID,
			Used: true,
		}

		if err := s.db.UpdateEmailVerificationToken(tx, updateEmailToken); err != nil {
			return status.Error(codes.Internal, "failed to update email verification token")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.VerifyEmailResponse{}, nil
}
