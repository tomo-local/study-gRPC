package service

import (
	"context"

	pb "auth/grpc/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerifyToken はJWTトークンの検証を行います
func (s *authServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	claims, err := s.jwtManager.VerifyToken(req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid: false,
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.VerifyTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
		Name:   claims.Name,
	}, nil
}
