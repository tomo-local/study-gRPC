package service

import (
	"context"
	"log"

	"note/db"
	"note/db/model"
	pb "note/grpc"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NoteServer struct {
	pb.UnimplementedNoteServiceServer
	db *db.DB
}

// NewNoteServer creates a new note server instance
func NewNoteServer(database *db.DB) *NoteServer {
	return &NoteServer{db: database}
}

// CreateNote creates a new note
func (s *NoteServer) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	note := &model.Note{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     pq.StringArray(req.Tags),
	}

	if err := s.db.Create(note).Error; err != nil {
		log.Printf("❌ Failed to create note: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create note: %v", err)
	}

	log.Printf("✅ Note created successfully with ID: %s", note.ID)
	return &pb.CreateNoteResponse{
		Note: modelToProto(note),
	}, nil
}

// GetNote retrieves a note by ID
func (s *NoteServer) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	log.Printf("🔍 Getting note: %s", req.Id)

	var note model.Note
	if err := s.db.First(&note, "id = ?", req.Id).Error; err != nil {
		log.Printf("❌ Failed to get note: %v", err)
		return nil, status.Errorf(codes.NotFound, "note not found: %v", err)
	}

	log.Printf("✅ Note retrieved successfully: %s", note.Title)
	return &pb.GetNoteResponse{
		Note: modelToProto(&note),
	}, nil
}

// ListNotes lists notes with pagination and filtering
func (s *NoteServer) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	log.Printf("📋 Listing notes (page: %d, limit: %d)", req.Page, req.Limit)

	var notes []model.Note
	var totalCount int64

	query := s.db.Model(&model.Note{})

	// Filter by category if provided
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// Filter by tags if provided
	if len(req.Tags) > 0 {
		query = query.Where("tags && ?", pq.Array(req.Tags))
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("❌ Failed to count notes: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to count notes: %v", err)
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(int(offset)).Limit(int(req.Limit)).Find(&notes).Error; err != nil {
		log.Printf("❌ Failed to list notes: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	protoNotes := make([]*pb.Note, len(notes))
	for i, note := range notes {
		protoNotes[i] = modelToProto(&note)
	}

	log.Printf("✅ Listed %d notes (total: %d)", len(notes), totalCount)
	return &pb.ListNotesResponse{
		Notes:      protoNotes,
		TotalCount: int32(totalCount),
	}, nil
}

// UpdateNote updates an existing note
func (s *NoteServer) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	log.Printf("✏️  Updating note: %s", req.Id)

	var note model.Note
	if err := s.db.First(&note, "id = ?", req.Id).Error; err != nil {
		log.Printf("❌ Failed to find note for update: %v", err)
		return nil, status.Errorf(codes.NotFound, "note not found: %v", err)
	}

	// Update fields
	note.Title = req.Title
	note.Content = req.Content
	note.Category = req.Category
	note.Tags = pq.StringArray(req.Tags)

	if err := s.db.Save(&note).Error; err != nil {
		log.Printf("❌ Failed to update note: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update note: %v", err)
	}

	log.Printf("✅ Note updated successfully: %s", note.Title)
	return &pb.UpdateNoteResponse{
		Note: modelToProto(&note),
	}, nil
}

// DeleteNote deletes a note by ID
func (s *NoteServer) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	log.Printf("🗑️  Deleting note: %s", req.Id)

	result := s.db.Delete(&model.Note{}, "id = ?", req.Id)
	if result.Error != nil {
		log.Printf("❌ Failed to delete note: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "failed to delete note: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("❌ Note not found for deletion: %s", req.Id)
		return nil, status.Errorf(codes.NotFound, "note not found")
	}

	log.Printf("✅ Note deleted successfully")
	return &pb.DeleteNoteResponse{
		Success: true,
	}, nil
}

// modelToProto converts a Note model to protobuf Note
func modelToProto(note *model.Note) *pb.Note {
	return &pb.Note{
		Id:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		Category:  note.Category,
		Tags:      []string(note.Tags),
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: timestamppb.New(note.UpdatedAt),
	}
}
