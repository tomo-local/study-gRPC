package model_test

import (
	"testing"
	"time"

	"email-service/db/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate the schema
	db.AutoMigrate(&model.User{})

	return db
}

func TestUser_Creation(t *testing.T) {
	db := setupTestDB()

	// Red: この時点ではUserモデルは存在しない
	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FullName:     "Test User",
		IsActive:     true,
	}

	err := db.Create(user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)
	assert.True(t, user.IsActive)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUser_UniqueEmail(t *testing.T) {
	db := setupTestDB()

	// Create first user
	user1 := &model.User{
		Email:        "unique@example.com",
		PasswordHash: "hash1",
		FullName:     "User One",
	}
	err := db.Create(user1).Error
	assert.NoError(t, err)

	// Try to create second user with same email
	user2 := &model.User{
		Email:        "unique@example.com", // Same email
		PasswordHash: "hash2",
		FullName:     "User Two",
	}
	err = db.Create(user2).Error
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestUser_Validation(t *testing.T) {
	// Test email validation (should not be empty)
	user := &model.User{
		Email:        "", // Empty email
		PasswordHash: "hash",
		FullName:     "Test User",
	}

	err := user.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email")
}

func TestUser_BeforeCreateHook(t *testing.T) {
	db := setupTestDB()

	user := &model.User{
		Email:        "hook@example.com",
		PasswordHash: "hash",
		FullName:     "Hook Test",
	}

	// CreatedAt should be zero before creation
	assert.True(t, user.CreatedAt.IsZero())

	err := db.Create(user).Error
	assert.NoError(t, err)

	// CreatedAt should be set after creation
	assert.False(t, user.CreatedAt.IsZero())
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
}
