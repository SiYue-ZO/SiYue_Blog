package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type contactRepo struct{ db *gorm.DB }

func NewContactRepo(db *gorm.DB) *contactRepo { return &contactRepo{db: db} }

func (r *contactRepo) Create(ctx context.Context, message *entity.ContactMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *contactRepo) List(ctx context.Context, offset, limit int, status string) ([]entity.ContactMessage, int64, error) {
	var messages []entity.ContactMessage
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.ContactMessage{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count contact messages: %w", err)
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("list contact messages: %w", err)
	}
	return messages, total, nil
}

func (r *contactRepo) GetByID(ctx context.Context, id uint64) (*entity.ContactMessage, error) {
	var message entity.ContactMessage
	if err := r.db.WithContext(ctx).First(&message, id).Error; err != nil {
		return nil, fmt.Errorf("get contact message: %w", err)
	}
	return &message, nil
}

func (r *contactRepo) UpdateStatus(ctx context.Context, id uint64, status entity.ContactStatus) error {
	result := r.db.WithContext(ctx).Model(&entity.ContactMessage{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *contactRepo) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&entity.ContactMessage{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
