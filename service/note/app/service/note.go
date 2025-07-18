package service

import (
	"context"
	"note/db"
	pb "note/grpc"
)

type NoteService interface {
	CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error)
	UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error)
	DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error)
	GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error)
	ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error)
}

type noteServer struct {
	pb.UnimplementedNoteServiceServer
	db db.DB
}

// NewNoteServer creates a new note server instance
func NewNoteServer(db db.DB) pb.NoteServiceServer {
	return &noteServer{db: db}
}
