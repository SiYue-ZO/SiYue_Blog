package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type homePageRepo struct {
	db *gorm.DB
}

type homePageRow struct {
	DraftConfig     []byte
	PublishedConfig []byte
	ConfigVersion   int
	UpdatedBy       *uint64
	PublishedAt     time.Time
	UpdatedAt       time.Time
}

func NewHomePageRepo(db *gorm.DB) *homePageRepo {
	return &homePageRepo{db: db}
}

func decodeHomePage(row homePageRow) (*entity.HomePageContent, error) {
	content := &entity.HomePageContent{
		ConfigVersion: row.ConfigVersion,
		UpdatedBy:     row.UpdatedBy,
		PublishedAt:   row.PublishedAt,
		UpdatedAt:     row.UpdatedAt,
	}
	if err := json.Unmarshal(row.DraftConfig, &content.DraftConfig); err != nil {
		return nil, fmt.Errorf("decode draft home page config: %w", err)
	}
	if err := json.Unmarshal(row.PublishedConfig, &content.PublishedConfig); err != nil {
		return nil, fmt.Errorf("decode published home page config: %w", err)
	}
	return content, nil
}

func (r *homePageRepo) Get(ctx context.Context) (*entity.HomePageContent, error) {
	var row homePageRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT draft_config, published_config, config_version, updated_by, published_at, updated_at
		FROM home_page_content WHERE id = 1
	`).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("get home page content: %w", err)
	}
	if len(row.DraftConfig) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return decodeHomePage(row)
}

func (r *homePageRepo) SaveDraft(ctx context.Context, config entity.HomePageConfig, userID uint64) (*entity.HomePageContent, error) {
	value, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("encode home page config: %w", err)
	}
	result := r.db.WithContext(ctx).Exec(`
		UPDATE home_page_content SET draft_config = ?::jsonb, updated_by = ? WHERE id = 1
	`, string(value), userID)
	if result.Error != nil {
		return nil, fmt.Errorf("save home page draft: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.Get(ctx)
}

func (r *homePageRepo) Publish(ctx context.Context, userID uint64) (*entity.HomePageContent, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE home_page_content
		SET published_config = draft_config, config_version = config_version + 1,
		    updated_by = ?, published_at = NOW()
		WHERE id = 1
	`, userID)
	if result.Error != nil {
		return nil, fmt.Errorf("publish home page: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.Get(ctx)
}

func (r *homePageRepo) ResetDraft(ctx context.Context, userID uint64) (*entity.HomePageContent, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE home_page_content SET draft_config = published_config, updated_by = ? WHERE id = 1
	`, userID)
	if result.Error != nil {
		return nil, fmt.Errorf("reset home page draft: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.Get(ctx)
}
