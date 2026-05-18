package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type tagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *tagRepo {
	return &tagRepo{db: db}
}

func (r *tagRepo) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepo) GetByID(ctx context.Context, id uint64) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, fmt.Errorf("get tag by id: %w", err)
	}
	return &tag, nil
}

func (r *tagRepo) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, fmt.Errorf("get tag by slug: %w", err)
	}
	return &tag, nil
}

func (r *tagRepo) List(ctx context.Context) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).
		Select("tags.*, COUNT(post_tags.post_id) as post_count").
		Joins("LEFT JOIN post_tags ON post_tags.tag_id = tags.id").
		Joins("LEFT JOIN posts ON posts.id = post_tags.post_id AND posts.status = ?", entity.PostStatusPublished).
		Group("tags.id").
		Order("post_count DESC").
		Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (r *tagRepo) ListByPostID(ctx context.Context, postID uint64) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("list tags by post id: %w", err)
	}
	return tags, nil
}

func (r *tagRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, id).Error
}
