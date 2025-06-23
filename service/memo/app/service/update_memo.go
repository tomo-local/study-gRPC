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

	updateMemo := &model.Memo{
		ID:       originMemo.ID,
		FileType: originMemo.FileType,
		Title:    originMemo.Title,
		Content:  req.Content,
	}

	updatedMemo, err := s.FileService.UpdateFile(updateMemo)
	if err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return &grpcPkg.UpdateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:        updatedMemo.ID,
			Title:     updatedMemo.Title,
			Content:   updatedMemo.Content,
			CreatedAt: timePkg.New(updatedMemo.CreatedAt),
			UpdatedAt: timePkg.New(updatedMemo.UpdatedAt),
		},
	}, nil
}
