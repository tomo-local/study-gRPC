package service

import (
	"context"
	"fmt"
	pb "note/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *noteServer) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	note, err := s.db.GetNoteByID(nil, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	return &pb.GetNoteResponse{
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
