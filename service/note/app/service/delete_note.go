package service

import (
	"context"
	"fmt"
	pb "note/grpc"
)

func (s *noteServer) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	// まず該当のノートが存在するか確認
	_, err := s.db.GetNoteByID(nil, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to find note for deletion: %w", err)
	}

	// ノートを削除
	err = s.db.DeleteNote(nil, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete note: %w", err)
	}

	return &pb.DeleteNoteResponse{
		Success: true,
	}, nil
}
