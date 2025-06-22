package service

import (
	pb "memo/grpc"
)

type MemoService struct {
	pb.UnimplementedMemoServiceServer
}

func NewMemoService() *MemoService {
	return &MemoService{}
}
