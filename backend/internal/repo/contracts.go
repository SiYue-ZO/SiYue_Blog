package repo

import (
	"context"
	"errors"

	"siyue-blog/internal/entity"
)

var ErrSessionReused = errors.New("refresh session reused")

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type SessionRepo interface {
	Create(ctx context.Context, session *entity.AuthSession) error
	Rotate(ctx context.Context, oldHash string, replacement *entity.AuthSession) error
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID uint64) error
}

type PostRepo interface {
	Create(ctx context.Context, post *entity.Post) error
	CreateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag) error
	GetByID(ctx context.Context, id uint64) (*entity.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Post, error)
	List(ctx context.Context, offset, limit int, tag string, search string) ([]entity.Post, int64, error)
	ListAdmin(ctx context.Context, offset, limit int, status string, search string) ([]entity.Post, int64, error)
	ListFeatured(ctx context.Context, limit int) ([]entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	UpdateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag, replaceTags bool) error
	Delete(ctx context.Context, id uint64) error
	IncrementViewCount(ctx context.Context, id uint64) error
}

type ContactRepo interface {
	Create(ctx context.Context, message *entity.ContactMessage) error
	List(ctx context.Context, offset, limit int, status string) ([]entity.ContactMessage, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.ContactMessage, error)
	UpdateStatus(ctx context.Context, id uint64, status entity.ContactStatus) error
	Delete(ctx context.Context, id uint64) error
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
	List(ctx context.Context, visibleOnly bool) ([]entity.Project, error)
	ListFeatured(ctx context.Context) ([]entity.Project, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id uint64) error
}

type SettingsRepo interface {
	Get(ctx context.Context) (*entity.SiteSettings, error)
	Update(ctx context.Context, projectPageVisible *bool, socialLinks *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error)
}

type HomePageRepo interface {
	Get(ctx context.Context) (*entity.HomePageContent, error)
	SaveDraft(ctx context.Context, config entity.HomePageConfig, userID uint64) (*entity.HomePageContent, error)
	Publish(ctx context.Context, userID uint64) (*entity.HomePageContent, error)
	ResetDraft(ctx context.Context, userID uint64) (*entity.HomePageContent, error)
}
