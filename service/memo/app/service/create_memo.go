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
	fileType := model.FileTypeTxt // デフォルトはtxt。拡張したい場合はリクエストに追加してね！
	now := timePkg.Now().AsTime()

	memo := &model.Memo{
		ID:        id,
		FileType:  fileType,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
	}

	if _, err := s.FileService.CreateFile(memo); err != nil {
		return nil, fmt.Errorf("ファイル保存失敗: %w", err)
	}

	return &grpcPkg.CreateMemoResponse{
		Memo: &grpcPkg.Memo{
			Id:        memo.ID,
			Title:     memo.Title,
			Content:   memo.Content,
			CreatedAt: timePkg.New(now),
			UpdatedAt: timePkg.New(now),
		},
	}, nil
}
