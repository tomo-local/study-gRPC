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
	// ID生成
	id := uuid.New().String()
	fileType := model.FileTypeMd
	now := timePkg.Now().AsTime()

	memo := &model.Memo{
		ID:         id,
		FileType:   fileType,
		Title:      req.Title,
		Content:    req.Content,
		ModifiedAt: now,
	}

	if _, err := s.FileService.CreateFile(memo); err != nil {
		return nil, fmt.Errorf("ファイル保存失敗: %w", err)
	}

	return &grpcPkg.CreateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:         memo.ID,
			Title:      memo.Title,
			Content:    memo.Content,
			ModifiedAt: timePkg.New(now),
		},
	}, nil
}
