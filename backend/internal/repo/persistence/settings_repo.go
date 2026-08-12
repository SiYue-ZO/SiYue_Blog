package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type settingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepo(db *gorm.DB) *settingsRepo {
	return &settingsRepo{db: db}
}

func (r *settingsRepo) Get(ctx context.Context) (*entity.SiteSettings, error) {
	var settings entity.SiteSettings
	if err := r.db.WithContext(ctx).First(&settings, 1).Error; err != nil {
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	return &settings, nil
}

func (r *settingsRepo) Update(ctx context.Context, projectPageVisible *bool, socialLinks *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error) {
	updates := entity.SiteSettings{ID: 1}
	fields := make([]string, 0, 3)
	if projectPageVisible != nil {
		updates.ProjectPageVisible = *projectPageVisible
		fields = append(fields, "project_page_visible")
	}
	if socialLinks != nil {
		updates.SocialLinks = socialLinks
		fields = append(fields, "social_links")
	}
	if footer != nil {
		updates.Footer = footer
		fields = append(fields, "footer")
	}
	result := r.db.WithContext(ctx).
		Model(&entity.SiteSettings{}).
		Where("id = ?", 1).
		Select(fields).
		Updates(&updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update site settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.Get(ctx)
}
