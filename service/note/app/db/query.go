package db

import (
	"fmt"
	"note/db/model"

	"gorm.io/gorm"
)

func (d *db) CreateNote(tx *gorm.DB, note *model.Note) error {
	client := d.getClient(tx)

	if err := client.Create(note).Error; err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

func (d *db) UpdateNote(tx *gorm.DB, note *model.Note) error {
	client := d.getClient(tx)

	if err := client.Save(note).Error; err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}

func (d *db) DeleteNote(tx *gorm.DB, id string) error {
	client := d.getClient(tx)

	if err := client.Delete(&model.Note{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete note with ID %s: %w", id, err)
	}

	return nil
}

func (d *db) GetNoteByID(tx *gorm.DB, id string) (*model.Note, error) {
	client := d.getClient(tx)

	var note model.Note
	if err := client.First(&note, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("note with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get note by ID %s: %w", id, err)
	}

	return &note, nil
}

func (d *db) ListNotes(tx *gorm.DB, page, limit int32, category string, tags []string) ([]*model.Note, int64, error) {
	client := d.getClient(tx)

	var notes []*model.Note
	var totalCount int64

	// デフォルト値の設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// クエリビルダーを作成
	query := client.Model(&model.Note{})

	// カテゴリフィルタ
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// タグフィルタ
	if len(tags) > 0 {
		query = query.Where("tags @> ?", tags)
	}

	// 総件数を取得
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count notes: %w", err)
	}

	// ページネーション付きでノートを取得
	if err := query.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list notes: %w", err)
	}

	return notes, totalCount, nil
}
