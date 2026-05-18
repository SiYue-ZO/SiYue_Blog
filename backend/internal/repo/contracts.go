package repo

import (
	"context"

	"siyue-blog/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type PostRepo interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id uint64) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, offset, limit int, tag string, search string) ([]entity.Post, int64, error)
	ListFeatured(ctx context.Context, limit int) ([]entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uint64) error
	IncrementViewCount(ctx context.Context, id uint64) error
}

type TagRepo interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uint64) (*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	List(ctx context.Context) ([]entity.Tag, error)
	ListByPostID(ctx context.Context, postID uint64) ([]entity.Tag, error)
	Delete(ctx context.Context, id uint64) error
}

type ProjectRepo interface {
	Create(ctx context.Context, project *entity.Project) error
	GetByID(ctx context.Context, id uint64) (*entity.Project, error)
	List(ctx context.Context) ([]entity.Project, error)
	ListFeatured(ctx context.Context) ([]entity.Project, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id uint64) error
}
