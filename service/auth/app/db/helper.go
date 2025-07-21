package db

import (
	"gorm.io/gorm"
)

func (d *db) StartTransaction(fn func(tx *gorm.DB) error) error {
	if err := d.client.Transaction(fn); err != nil {
		return err
	}
	return nil
}

func (d *db) getClient(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return d.client
	}
	return tx
}
