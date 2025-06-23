package service

import (
	"context"
	"fmt"
	grpcPkg "memo/grpc"

	timePkg "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MemoService) GetMemo(ctx context.Context, req *grpcPkg.GetMemoRequest) (*grpcPkg.GetMemoResponse, error) {
	memo, err := s.FileService.GetFile(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memo %w", err)
	}

	return &grpcPkg.GetMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:         memo.ID,
			Title:      memo.Title,
			Content:    memo.Content,
			ModifiedAt: timePkg.New(memo.ModifiedAt),
		},
	}, nil

}
