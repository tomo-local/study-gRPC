package service

import (
	"context"
	"fmt"
	"memo/db/model"
	grpcPkg "memo/grpc"

	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) UpdateMemo(ctx context.Context, req *grpcPkg.UpdateMemoRequest) (*grpcPkg.UpdateMemoResponse, error) {
	originMemo, err := s.FileService.GetFile(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memo %w", err)
	}

	now := timePkg.Now().AsTime()

	updateMemo := &model.Memo{
		ID:         originMemo.ID,
		Title:      req.Title,
		Content:    req.Content,
		ModifiedAt: now,
	}

	if _, err := s.FileService.UpdateFile(updateMemo); err != nil {
		return nil, fmt.Errorf("failed to update memo %w", err)
	}

	return &grpcPkg.UpdateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:         updateMemo.ID,
			Title:      updateMemo.Title,
			Content:    updateMemo.Content,
			ModifiedAt: timePkg.New(now),
		},
	}, nil
}
