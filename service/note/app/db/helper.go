package db

import (
	"fmt"

	"gorm.io/gorm"
)

func (d *db) StartTransaction(fn func(tx *gorm.DB) error) error {
	if err := d.client.Transaction(fn); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (d *db) getClient(tx *gorm.DB) *gorm.DB {
	client := d.client
	if tx != nil {
		client = tx
	}
	return client
}
