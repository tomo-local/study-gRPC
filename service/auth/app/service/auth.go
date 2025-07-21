package service

import (
	"context"

	"auth/db"
	pb "auth/grpc/api"
	"auth/jwt"
	"auth/mailer"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error)
	RequestPasswordReset(ctx context.Context, req *pb.RequestPasswordResetRequest) (*pb.RequestPasswordResetResponse, error)
	ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error)
	ResendVerificationEmail(ctx context.Context, req *pb.ResendVerificationEmailRequest) (*pb.ResendVerificationEmailResponse, error)
}

type authServer struct {
	pb.UnimplementedAuthServiceServer
	db         db.DB
	jwtManager *jwt.Manager
	mailer     *mailer.Mailer
}

// NewAuthServer creates a new auth server instance
func NewAuthServer(database db.DB, jwtManager *jwt.Manager, mailer *mailer.Mailer) pb.AuthServiceServer {
	return &authServer{
		db:         database,
		jwtManager: jwtManager,
		mailer:     mailer,
	}
}
