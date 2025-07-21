package service

import (
	"context"
	"fmt"

	"auth/db/model"
	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Login はユーザーのログインを行います
func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user *model.User
	var accessToken, refreshToken string

	err := s.db.StartTransaction(func(tx *gorm.DB) error {
		// ユーザーを取得
		var err error
		user, err = s.db.GetUserByEmail(tx, req.Email)
		if err != nil {
			return status.Error(codes.Unauthenticated, "invalid email or password")
		}

		// パスワードの検証
		if !checkPassword(req.Password, user.PasswordHash) {
			return status.Error(codes.Unauthenticated, "invalid credentials")
		} // メールアドレスが認証済みかチェック
		if !user.EmailVerified {
			return status.Error(codes.Unauthenticated, "email not verified")
		}

		// アクセストークンを生成
		accessToken, err = s.jwtManager.GenerateToken(user.ID, user.Email, user.Name)
		if err != nil {
			return fmt.Errorf("failed to generate access token: %w", err)
		}

		// リフレッシュトークン（実際の実装では別のトークン生成方法を使用）
		refreshToken, err = s.jwtManager.GenerateToken(user.ID, user.Email, user.Name)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         convertUserToProto(user),
	}, nil
}
