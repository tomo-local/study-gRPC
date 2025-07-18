package db

import (
	"auth/db/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User関連のメソッド

// CreateUser は新しいユーザーを作成します。
func (d *db) CreateUser(tx *gorm.DB, user *model.User) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID はIDからユーザーを取得します。
func (d *db) GetUserByID(tx *gorm.DB, id string) (*model.User, error) {
	if tx == nil {
		tx = d.client
	}

	var user model.User
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetUserByEmail はメールアドレスからユーザーを取得します。
func (d *db) GetUserByEmail(tx *gorm.DB, email string) (*model.User, error) {
	if tx == nil {
		tx = d.client
	}

	var user model.User
	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// UpdateUser はユーザー情報を更新します。
func (d *db) UpdateUser(tx *gorm.DB, user *model.User) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser はユーザーを削除します。
func (d *db) DeleteUser(tx *gorm.DB, id string) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// EmailVerificationToken関連のメソッド

// CreateEmailVerificationToken は新しいメールアドレス認証トークンを作成します。
func (d *db) CreateEmailVerificationToken(tx *gorm.DB, token *model.EmailVerificationToken) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Create(token).Error; err != nil {
		return fmt.Errorf("failed to create email verification token: %w", err)
	}

	return nil
}

// GetEmailVerificationTokenByToken はトークンからメールアドレス認証トークンを取得します。
func (d *db) GetEmailVerificationTokenByToken(tx *gorm.DB, token string) (*model.EmailVerificationToken, error) {
	if tx == nil {
		tx = d.client
	}

	var emailToken model.EmailVerificationToken
	if err := tx.Where("token = ?", token).First(&emailToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email verification token not found: %s", token)
		}
		return nil, fmt.Errorf("failed to get email verification token: %w", err)
	}

	return &emailToken, nil
}

// UpdateEmailVerificationToken はメールアドレス認証トークンを更新します。
func (d *db) UpdateEmailVerificationToken(tx *gorm.DB, token *model.EmailVerificationToken) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Save(token).Error; err != nil {
		return fmt.Errorf("failed to update email verification token: %w", err)
	}

	return nil
}

// DeleteEmailVerificationToken はメールアドレス認証トークンを削除します。
func (d *db) DeleteEmailVerificationToken(tx *gorm.DB, id string) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Delete(&model.EmailVerificationToken{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete email verification token: %w", err)
	}

	return nil
}

// PasswordResetToken関連のメソッド

// CreatePasswordResetToken は新しいパスワードリセットトークンを作成します。
func (d *db) CreatePasswordResetToken(tx *gorm.DB, token *model.PasswordResetToken) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Create(token).Error; err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	return nil
}

// GetPasswordResetTokenByToken はトークンからパスワードリセットトークンを取得します。
func (d *db) GetPasswordResetTokenByToken(tx *gorm.DB, token string) (*model.PasswordResetToken, error) {
	if tx == nil {
		tx = d.client
	}

	var resetToken model.PasswordResetToken
	if err := tx.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("password reset token not found: %s", token)
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}

	return &resetToken, nil
}

// UpdatePasswordResetToken はパスワードリセットトークンを更新します。
func (d *db) UpdatePasswordResetToken(tx *gorm.DB, token *model.PasswordResetToken) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Save(token).Error; err != nil {
		return fmt.Errorf("failed to update password reset token: %w", err)
	}

	return nil
}

// DeletePasswordResetToken はパスワードリセットトークンを削除します。
func (d *db) DeletePasswordResetToken(tx *gorm.DB, id string) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Delete(&model.PasswordResetToken{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}

	return nil
}

// DeleteExpiredPasswordResetTokens は期限切れのパスワードリセットトークンを削除します。
func (d *db) DeleteExpiredPasswordResetTokens(tx *gorm.DB) error {
	if tx == nil {
		tx = d.client
	}

	if err := tx.Where("expires_at < ?", time.Now()).Delete(&model.PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete expired password reset tokens: %w", err)
	}

	return nil
}
