package service

import (
	"context"
	"fmt"
	"memo/db/model"
	grpcPkg "memo/grpc"
	"time"

	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) UpdateMemo(ctx context.Context, req *grpcPkg.UpdateMemoRequest) (*grpcPkg.UpdateMemoResponse, error) {
	originMemo, err := s.FileService.GetFile(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memo %w", err)
	}

	now := timePkg.Now().AsTime()

	updateMemo := &model.Memo{
		ID:        originMemo.ID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: originMemo.CreatedAt,
		UpdatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
	}

	if _, err := s.FileService.UpdateFile(updateMemo); err != nil {
		return nil, fmt.Errorf("failed to update memo %w", err)
	}

	// Parse memo.CreatedAt from string to time.Time
	createdAt, err := time.Parse("2006-01-02T15:04 :05Z07:00", originMemo.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to parse CreatedAt: %w", err)
	}

	return &grpcPkg.UpdateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:        updateMemo.ID,
			Title:     updateMemo.Title,
			Content:   updateMemo.Content,
			CreatedAt: timePkg.New(createdAt),
			UpdatedAt: timePkg.New(now),
		},
	}, nil
}
