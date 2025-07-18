package service

import (
	"context"
	"fmt"
	pb "note/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *noteServer) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	// デフォルト値の設定
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 10
	}

	// データベースからノートを取得
	notes, totalCount, err := s.db.ListNotes(nil, page, limit, req.GetCategory(), req.GetTags())
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	// gRPCレスポンス用に変換
	grpcNotes := make([]*pb.Note, len(notes))
	for i, note := range notes {
		grpcNotes[i] = &pb.Note{
			Id:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			Category:  note.Category,
			Tags:      note.Tags,
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: timestamppb.New(note.UpdatedAt),
		}
	}

	return &pb.ListNotesResponse{
		Notes:      grpcNotes,
		TotalCount: int32(totalCount),
	}, nil
}
