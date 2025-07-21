package service

import (
	"context"
	"time"

	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// VerifyEmail はメールアドレス認証を行います
func (s *authServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	err := s.db.StartTransaction(func(tx *gorm.DB) error {
		// トークンを取得
		emailToken, err := s.db.GetEmailVerificationTokenByToken(tx, req.Token)
		if err != nil {
			return err
		}

		// トークンが期限切れでないかチェック
		if emailToken.IsExpired() {
			return status.Error(codes.Unauthenticated, "verification token has expired")
		}

		// トークンが使用済みでないかチェック
		if emailToken.Used {
			return status.Error(codes.Unauthenticated, "verification token has already been used")
		}

		// ユーザーを取得
		user, err := s.db.GetUserByID(tx, emailToken.UserID)
		if err != nil {
			return err
		}

		// ユーザーのメールアドレスを認証済みに更新
		user.EmailVerified = true
		user.UpdatedAt = time.Now()

		if err := s.db.UpdateUser(tx, user); err != nil {
			return err
		}

		// トークンを使用済みにマーク
		emailToken.Used = true
		if err := s.db.UpdateEmailVerificationToken(tx, emailToken); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.VerifyEmailResponse{}, nil
}
