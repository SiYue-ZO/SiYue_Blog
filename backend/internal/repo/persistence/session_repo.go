package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
)

type sessionRepo struct{ db *gorm.DB }

func NewSessionRepo(db *gorm.DB) *sessionRepo { return &sessionRepo{db: db} }

func (r *sessionRepo) Create(ctx context.Context, session *entity.AuthSession) error {
	_ = r.db.WithContext(ctx).Where("expires_at < ?", time.Now().Add(-24*time.Hour)).Delete(&entity.AuthSession{}).Error
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepo) Rotate(ctx context.Context, oldHash string, replacement *entity.AuthSession) error {
	var outcomeErr error
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current entity.AuthSession
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", oldHash).First(&current).Error; err != nil {
			return fmt.Errorf("find session: %w", err)
		}
		if time.Now().After(current.ExpiresAt) || current.UserID != replacement.UserID {
			return gorm.ErrRecordNotFound
		}
		if current.RevokedAt != nil {
			now := time.Now()
			if err := tx.Model(&entity.AuthSession{}).
				Where("family_id = ? AND revoked_at IS NULL", current.FamilyID).
				Update("revoked_at", now).Error; err != nil {
				return fmt.Errorf("revoke reused session family: %w", err)
			}
			outcomeErr = repo.ErrSessionReused
			return nil
		}
		now := time.Now()
		if err := tx.Model(&current).Updates(map[string]any{"revoked_at": now, "replaced_by": replacement.TokenHash}).Error; err != nil {
			return err
		}
		replacement.FamilyID = current.FamilyID
		return tx.Create(replacement).Error
	})
	if err != nil {
		return err
	}
	return outcomeErr
}

func (r *sessionRepo) Revoke(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).Model(&entity.AuthSession{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", time.Now()).Error
}

func (r *sessionRepo) RevokeAllForUser(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Model(&entity.AuthSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now()).Error
}
