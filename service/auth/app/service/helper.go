package service

import (
	"time"

	"auth/db/model"
	pb "auth/grpc/api"
)

// convertUserToProto はmodel.UserをProtoのUserに変換します
func convertUserToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
	}
}
