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

func (r *postRepo) CreateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		if len(tags) > 0 {
			if err := tx.Model(post).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		post.Tags = tags
		return nil
	})
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
		Where("slug = ? AND status = ? AND published_at IS NOT NULL AND published_at <= NOW()", slug, entity.PostStatusPublished).
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
	query = query.Where("status = ? AND published_at IS NOT NULL AND published_at <= NOW()", entity.PostStatusPublished)

	// Filter by tag
	if tag != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.slug = ?", tag)
	}

	// Search filter
	if search != "" {
		query = query.Where("(posts.title || ' ' || COALESCE(posts.excerpt, '')) ILIKE ?", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	// Fetch posts with relations
	err := query.
		Omit("content").
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

func (r *postRepo) ListAdmin(ctx context.Context, offset, limit int, status string, search string) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.Post{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR slug ILIKE ? OR excerpt ILIKE ?", pattern, pattern, pattern)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count admin posts: %w", err)
	}
	if err := query.Omit("content").Preload("Tags").Preload("Author").
		Order("updated_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("list admin posts: %w", err)
	}
	return posts, total, nil
}

func (r *postRepo) ListFeatured(ctx context.Context, limit int) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.WithContext(ctx).
		Omit("content").
		Preload("Tags").
		Preload("Author").
		Where("status = ? AND is_pinned = ? AND published_at IS NOT NULL AND published_at <= NOW()", entity.PostStatusPublished, true).
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("list featured posts: %w", err)
	}
	return posts, nil
}

func (r *postRepo) Update(ctx context.Context, post *entity.Post) error {
	return r.UpdateWithTags(ctx, post, nil, false)
}

func (r *postRepo) UpdateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag, replaceTags bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.Post{}).Where("id = ?", post.ID).Updates(map[string]interface{}{
			"title": post.Title, "slug": post.Slug, "excerpt": post.Excerpt,
			"content": post.Content, "cover_image": post.CoverImage,
			"status": post.Status, "is_pinned": post.IsPinned, "published_at": post.PublishedAt,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if replaceTags {
			if err := tx.Model(post).Association("Tags").Replace(tags); err != nil {
				return err
			}
			post.Tags = tags
		}
		return nil
	})
}

func (r *postRepo) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&entity.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *postRepo) RecordView(ctx context.Context, id uint64, eventID string) (int, error) {
	var viewCount int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock first so duplicate events and different readers see a consistent
		// count, and an unpublish/delete cannot race the increment.
		result := tx.Raw(`SELECT view_count FROM posts
			WHERE id = ? AND status = 'published' AND published_at <= NOW()
			FOR UPDATE`, id).Scan(&viewCount)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		result = tx.Exec(`INSERT INTO post_view_events (post_id, event_id)
			VALUES (?, ?) ON CONFLICT DO NOTHING`, id, eventID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		return tx.Raw("UPDATE posts SET view_count = view_count + 1 WHERE id = ? RETURNING view_count", id).
			Scan(&viewCount).Error
	})
	return viewCount, err
}
