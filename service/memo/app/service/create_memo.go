package service

import (
	"context"
	"fmt"
	"memo/db/model"
	grpcPkg "memo/grpc"

	"github.com/google/uuid"
	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) CreateMemo(ctx context.Context, req *grpcPkg.CreateMemoRequest) (*grpcPkg.CreateMemoResponse, error) {
	memo := &model.Memo{
		ID:       uuid.New().String(),
		FileType: model.FileTypeMd,
		Title:    req.Title,
		Content:  req.Content,
	}

	createdMemo, err := s.FileService.CreateFile(memo)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return &grpcPkg.CreateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:        createdMemo.ID,
			Title:     createdMemo.Title,
			Content:   createdMemo.Content,
			CreatedAt: timePkg.New(createdMemo.CreatedAt),
			UpdatedAt: timePkg.New(createdMemo.UpdatedAt),
		},
	}, nil
}
