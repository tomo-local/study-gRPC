package service

import (
	"context"
	"fmt"
	grpcPkg "memo/grpc"

	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) ListMemos(ctx context.Context, req *grpcPkg.ListMemosRequest) (*grpcPkg.ListMemosResponse, error) {
	memos, err := s.FileService.ListFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list memos: %w", err)
	}

	grpcMemos := make([]*grpcPkg.Memo, 0)
	for _, memo := range memos {
		grpcMemos = append(grpcMemos, &grpcPkg.Memo{
			Id:        memo.ID,
			Title:     memo.Title,
			Content:   memo.Content,
			CreatedAt: timePkg.New(memo.CreatedAt),
			UpdatedAt: timePkg.New(memo.UpdatedAt),
		})
	}

	return &grpcPkg.ListMemosResponse{Memos: grpcMemos}, nil
}
