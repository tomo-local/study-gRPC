package service

import (
	"context"
	"fmt"
	grpcPkg "memo/grpc"
	"time"

	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) GetMemo(ctx context.Context, req *grpcPkg.GetMemoRequest) (*grpcPkg.GetMemoResponse, error) {
	memo, err := s.FileService.GetFile(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memo %w", err)
	}
	// Parse memo.CreatedAt and memo.UpdatedAt from string to time.Time
	createdAt, err := time.Parse("2006-01-02 15:04:05", memo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CreatedAt: %w", err)
	}
	updatedAt, err := time.Parse("2006-01-02 15:04:05", memo.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UpdatedAt: %w", err)
	}

	return &grpcPkg.GetMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:        memo.ID,
			Title:     memo.Title,
			Content:   memo.Content,
			CreatedAt: timePkg.New(createdAt),
			UpdatedAt: timePkg.New(updatedAt),
		},
	}, nil

}
