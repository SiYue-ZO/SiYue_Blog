package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *projectRepo {
	return &projectRepo{db: db}
}

func (r *projectRepo) Create(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepo) GetByID(ctx context.Context, id uint64) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).First(&project, id).Error
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	return &project, nil
}

func (r *projectRepo) List(ctx context.Context, visibleOnly bool) ([]entity.Project, error) {
	var projects []entity.Project
	query := r.db.WithContext(ctx)
	if visibleOnly {
		query = query.Where("is_visible = ?", true)
	}
	err := query.
		Order("sort_order ASC, created_at DESC").
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

func (r *projectRepo) ListFeatured(ctx context.Context) ([]entity.Project, error) {
	var projects []entity.Project
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND is_visible = ?", true, true).
		Order("sort_order ASC").
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("list featured projects: %w", err)
	}
	return projects, nil
}

func (r *projectRepo) Update(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepo) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&entity.Project{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
