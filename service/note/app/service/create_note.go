package service

import (
	"context"
	"fmt"
	"note/db/model"
	pb "note/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *noteServer) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	note := &model.Note{
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		Category: req.GetCategory(),
		Tags:     req.GetTags(),
	}

	err := s.db.CreateNote(nil, note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return &pb.CreateNoteResponse{
		Note: &pb.Note{
			Id:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			Category:  note.Category,
			Tags:      note.Tags,
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: timestamppb.New(note.UpdatedAt),
		},
	}, nil
}
