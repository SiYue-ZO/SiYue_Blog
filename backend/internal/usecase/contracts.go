package usecase

import (
	"context"

	"siyue-blog/internal/entity"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *entity.User, err error)
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type PostUseCase interface {
	Create(ctx context.Context, input CreatePostInput) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, params ListPostsParams) ([]entity.Post, int64, error)
	ListFeatured(ctx context.Context, limit int) ([]entity.Post, error)
	Update(ctx context.Context, id uint64, input UpdatePostInput) (*entity.Post, error)
	Delete(ctx context.Context, id uint64) error
	IncrementViewCount(ctx context.Context, id uint64) error
}

type TagUseCase interface {
	List(ctx context.Context) ([]entity.Tag, error)
	Create(ctx context.Context, name, slug string) (*entity.Tag, error)
	Delete(ctx context.Context, id uint64) error
}

type ProjectUseCase interface {
	List(ctx context.Context) ([]entity.Project, error)
	ListFeatured(ctx context.Context) ([]entity.Project, error)
	Create(ctx context.Context, input CreateProjectInput) (*entity.Project, error)
	Update(ctx context.Context, id uint64, input UpdateProjectInput) (*entity.Project, error)
	Delete(ctx context.Context, id uint64) error
}

type CreatePostInput struct {
	Title      string   `json:"title" binding:"required"`
	Slug       string   `json:"slug" binding:"required"`
	Excerpt    string   `json:"excerpt"`
	Content    string   `json:"content"`
	CoverImage string   `json:"cover_image"`
	Status     string   `json:"status"`
	IsPinned   bool     `json:"is_pinned"`
	TagIDs     []uint64 `json:"tag_ids"`
}

type UpdatePostInput struct {
	Title      *string  `json:"title"`
	Slug       *string  `json:"slug"`
	Excerpt    *string  `json:"excerpt"`
	Content    *string  `json:"content"`
	CoverImage *string  `json:"cover_image"`
	Status     *string  `json:"status"`
	IsPinned   *bool    `json:"is_pinned"`
	TagIDs     []uint64 `json:"tag_ids"`
}

type ListPostsParams struct {
	Page     int
	PageSize int
	Tag      string
	Search   string
}

type CreateProjectInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	RepoURL     string   `json:"repo_url"`
	DemoURL     string   `json:"demo_url"`
	ImageURL    string   `json:"image_url"`
	TechStack   []string `json:"tech_stack"`
	IsFeatured  bool     `json:"is_featured"`
	SortOrder   int      `json:"sort_order"`
}

type UpdateProjectInput struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	RepoURL     *string  `json:"repo_url"`
	DemoURL     *string  `json:"demo_url"`
	ImageURL    *string  `json:"image_url"`
	TechStack   []string `json:"tech_stack"`
	IsFeatured  *bool    `json:"is_featured"`
	SortOrder   *int     `json:"sort_order"`
}
