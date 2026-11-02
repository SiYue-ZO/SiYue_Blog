package usecase

import (
	"context"
	"errors"

	"siyue-blog/internal/entity"
)

var ErrInvalidPasswordChange = errors.New("invalid password change")

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *entity.User, err error)
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
	ChangePassword(ctx context.Context, userID uint64, currentPassword, newPassword string) error
}

type PostUseCase interface {
	Create(ctx context.Context, input CreatePostInput) (*entity.Post, error)
	GetByID(ctx context.Context, id uint64) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, params ListPostsParams) ([]entity.Post, int64, error)
	ListAdmin(ctx context.Context, params ListAdminPostsParams) ([]entity.Post, int64, error)
	ListFeatured(ctx context.Context, limit int) ([]entity.Post, error)
	Update(ctx context.Context, id uint64, input UpdatePostInput) (*entity.Post, error)
	Delete(ctx context.Context, id uint64) error
	IncrementViewCount(ctx context.Context, id uint64) (int, error)
}

type ContactUseCase interface {
	Create(ctx context.Context, input CreateContactInput) (*entity.ContactMessage, error)
	List(ctx context.Context, page, pageSize int, status string) ([]entity.ContactMessage, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
	Delete(ctx context.Context, id uint64) error
}

type TagUseCase interface {
	List(ctx context.Context) ([]entity.Tag, error)
	Create(ctx context.Context, name, slug string) (*entity.Tag, error)
	Delete(ctx context.Context, id uint64) error
}

type ProjectUseCase interface {
	List(ctx context.Context) ([]entity.Project, error)
	ListAdmin(ctx context.Context) ([]entity.Project, error)
	ListFeatured(ctx context.Context) ([]entity.Project, error)
	Create(ctx context.Context, input CreateProjectInput) (*entity.Project, error)
	Update(ctx context.Context, id uint64, input UpdateProjectInput) (*entity.Project, error)
	Delete(ctx context.Context, id uint64) error
}

type SettingsUseCase interface {
	Get(ctx context.Context) (*entity.SiteSettings, error)
	Update(ctx context.Context, projectPageVisible *bool, socialLinks *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error)
}

type HomePageUseCase interface {
	GetPublished(ctx context.Context) (*entity.HomePageConfig, error)
	GetAdmin(ctx context.Context) (*entity.HomePageContent, error)
	SaveDraft(ctx context.Context, config entity.HomePageConfig, userID uint64) (*entity.HomePageContent, error)
	Publish(ctx context.Context, userID uint64) (*entity.HomePageContent, error)
	ResetDraft(ctx context.Context, userID uint64) (*entity.HomePageContent, error)
}

type CreatePostInput struct {
	Title      string   `json:"title" binding:"required,max=255"`
	Slug       string   `json:"slug" binding:"required,max=255"`
	Excerpt    string   `json:"excerpt" binding:"max=500"`
	Content    string   `json:"content"`
	CoverImage string   `json:"cover_image"`
	Status     string   `json:"status"`
	IsPinned   bool     `json:"is_pinned"`
	TagIDs     []uint64 `json:"tag_ids"`
	AuthorID   uint64   `json:"-"`
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

type ListAdminPostsParams struct {
	Page     int
	PageSize int
	Status   string
	Search   string
}

type CreateContactInput struct {
	Name    string
	Email   string
	Message string
	IPHash  string
}

type CreateProjectInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	RepoURL     string   `json:"repo_url"`
	DemoURL     string   `json:"demo_url"`
	ImageURL    string   `json:"image_url"`
	TechStack   []string `json:"tech_stack"`
	IsFeatured  bool     `json:"is_featured"`
	IsVisible   *bool    `json:"is_visible"`
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
	IsVisible   *bool    `json:"is_visible"`
	SortOrder   *int     `json:"sort_order"`
}
