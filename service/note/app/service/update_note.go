package service

import (
	"context"
	"fmt"
	"note/db/model"
	pb "note/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *noteServer) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	// 既存のノートを取得
	existingNote, err := s.db.GetNoteByID(nil, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get note for update: %w", err)
	}

	// 更新するフィールドを設定
	updatedNote := &model.Note{
		ID:       existingNote.ID,
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		Category: req.GetCategory(),
		Tags:     req.GetTags(),
		// CreatedAtは既存の値を保持
		CreatedAt: existingNote.CreatedAt,
		// UpdatedAtは自動更新される
	}

	// データベースで更新
	err = s.db.UpdateNote(nil, updatedNote)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	// 更新後のノートを取得（UpdatedAtを正確に取得するため）
	refreshedNote, err := s.db.GetNoteByID(nil, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get updated note: %w", err)
	}

	return &pb.UpdateNoteResponse{
		Note: &pb.Note{
			Id:        refreshedNote.ID,
			Title:     refreshedNote.Title,
			Content:   refreshedNote.Content,
			Category:  refreshedNote.Category,
			Tags:      refreshedNote.Tags,
			CreatedAt: timestamppb.New(refreshedNote.CreatedAt),
			UpdatedAt: timestamppb.New(refreshedNote.UpdatedAt),
		},
	}, nil
}
