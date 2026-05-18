package project

import (
	"context"
	"fmt"

	"github.com/lib/pq"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	projectRepo repo.ProjectRepo
}

func NewUseCase(projectRepo repo.ProjectRepo) usecase.ProjectUseCase {
	return &useCase{
		projectRepo: projectRepo,
	}
}

func (uc *useCase) List(ctx context.Context) ([]entity.Project, error) {
	projects, err := uc.projectRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

func (uc *useCase) ListFeatured(ctx context.Context) ([]entity.Project, error) {
	projects, err := uc.projectRepo.ListFeatured(ctx)
	if err != nil {
		return nil, fmt.Errorf("list featured projects: %w", err)
	}
	return projects, nil
}

func (uc *useCase) Create(ctx context.Context, input usecase.CreateProjectInput) (*entity.Project, error) {
	project := &entity.Project{
		Title:       input.Title,
		Description: input.Description,
		RepoURL:     input.RepoURL,
		DemoURL:     input.DemoURL,
		ImageURL:    input.ImageURL,
		TechStack:   pq.StringArray(input.TechStack),
		IsFeatured:  input.IsFeatured,
		SortOrder:   input.SortOrder,
	}

	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return project, nil
}

func (uc *useCase) Update(ctx context.Context, id uint64, input usecase.UpdateProjectInput) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	if input.Title != nil {
		project.Title = *input.Title
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.RepoURL != nil {
		project.RepoURL = *input.RepoURL
	}
	if input.DemoURL != nil {
		project.DemoURL = *input.DemoURL
	}
	if input.ImageURL != nil {
		project.ImageURL = *input.ImageURL
	}
	if input.TechStack != nil {
		project.TechStack = pq.StringArray(input.TechStack)
	}
	if input.IsFeatured != nil {
		project.IsFeatured = *input.IsFeatured
	}
	if input.SortOrder != nil {
		project.SortOrder = *input.SortOrder
	}

	if err := uc.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	return project, nil
}

func (uc *useCase) Delete(ctx context.Context, id uint64) error {
	if err := uc.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
