package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type postRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) Create(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepo) GetByID(ctx context.Context, id uint64) (*entity.Post, error) {
	var post entity.Post
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Author").
		First(&post, id).Error
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	return &post, nil
}

func (r *postRepo) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	var post entity.Post
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Author").
		Where("slug = ?", slug).
		First(&post).Error
	if err != nil {
		return nil, fmt.Errorf("get post by slug: %w", err)
	}
	return &post, nil
}

func (r *postRepo) List(ctx context.Context, offset, limit int, tag string, search string) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Post{})

	// Filter by status
	query = query.Where("status = ?", entity.PostStatusPublished)

	// Filter by tag
	if tag != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.slug = ?", tag)
	}

	// Search filter
	if search != "" {
		query = query.Where("posts.title ILIKE ? OR posts.excerpt ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	// Fetch posts with relations
	err := query.
		Preload("Tags").
		Preload("Author").
		Order("is_pinned DESC, published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}

	return posts, total, nil
}

func (r *postRepo) ListFeatured(ctx context.Context, limit int) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Author").
		Where("status = ? AND is_pinned = ?", entity.PostStatusPublished, true).
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("list featured posts: %w", err)
	}
	return posts, nil
}

func (r *postRepo) Update(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Post{}, id).Error
}

func (r *postRepo) IncrementViewCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&entity.Post{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
