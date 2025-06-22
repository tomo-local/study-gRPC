package service

import (
	config "memo/config/server"
	"memo/db"
	pb "memo/grpc"
)

type MemoService struct {
	pb.UnimplementedMemoServiceServer
	db.FileService
}

func NewMemoService(env *config.Config) *MemoService {
	fs, _ := db.GetService(env)
	return &MemoService{
		FileService: fs,
	}
}
