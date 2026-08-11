package project

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/lib/pq"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	projectRepo repo.ProjectRepo
}

func validateURL(value string) bool {
	if value == "" {
		return true
	}
	parsed, err := url.ParseRequestURI(value)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func normalizeStack(values []string) ([]string, error) {
	if len(values) > 20 {
		return nil, fmt.Errorf("tech stack must contain at most 20 items")
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || utf8.RuneCountInString(value) > 50 {
			return nil, fmt.Errorf("tech stack items must contain between 1 and 50 characters")
		}
		key := strings.ToLower(value)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
	}
	return result, nil
}

func validateProject(project *entity.Project) error {
	if project.Title == "" || utf8.RuneCountInString(project.Title) > 255 {
		return fmt.Errorf("title must contain between 1 and 255 characters")
	}
	for label, value := range map[string]string{
		"repository": project.RepoURL, "demo": project.DemoURL,
	} {
		if utf8.RuneCountInString(value) > 512 || !validateURL(value) {
			return fmt.Errorf("%s URL must be a valid HTTP(S) URL of at most 512 characters", label)
		}
	}
	if utf8.RuneCountInString(project.ImageURL) > 512 {
		return fmt.Errorf("image URL must not exceed 512 characters")
	}
	if project.ImageURL != "" {
		parsed, err := url.ParseRequestURI(project.ImageURL)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return fmt.Errorf("image URL must be a valid HTTPS URL")
		}
	}
	return nil
}

func NewUseCase(projectRepo repo.ProjectRepo) usecase.ProjectUseCase {
	return &useCase{
		projectRepo: projectRepo,
	}
}

func (uc *useCase) List(ctx context.Context) ([]entity.Project, error) {
	projects, err := uc.projectRepo.List(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

func (uc *useCase) ListAdmin(ctx context.Context) ([]entity.Project, error) {
	projects, err := uc.projectRepo.List(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("list admin projects: %w", err)
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
	input.Title = strings.TrimSpace(input.Title)
	isVisible := true
	if input.IsVisible != nil {
		isVisible = *input.IsVisible
	}
	stack, err := normalizeStack(input.TechStack)
	if err != nil {
		return nil, err
	}
	project := &entity.Project{
		Title:       input.Title,
		Description: input.Description,
		RepoURL:     input.RepoURL,
		DemoURL:     input.DemoURL,
		ImageURL:    input.ImageURL,
		TechStack:   pq.StringArray(stack),
		IsFeatured:  input.IsFeatured,
		IsVisible:   isVisible,
		SortOrder:   input.SortOrder,
	}
	if err := validateProject(project); err != nil {
		return nil, err
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
		project.Title = strings.TrimSpace(*input.Title)
		if project.Title == "" {
			return nil, fmt.Errorf("title is required")
		}
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.RepoURL != nil {
		if !validateURL(*input.RepoURL) {
			return nil, fmt.Errorf("invalid repository URL")
		}
		project.RepoURL = *input.RepoURL
	}
	if input.DemoURL != nil {
		if !validateURL(*input.DemoURL) {
			return nil, fmt.Errorf("invalid demo URL")
		}
		project.DemoURL = *input.DemoURL
	}
	if input.ImageURL != nil {
		if !validateURL(*input.ImageURL) {
			return nil, fmt.Errorf("invalid image URL")
		}
		project.ImageURL = *input.ImageURL
	}
	if input.TechStack != nil {
		stack, err := normalizeStack(input.TechStack)
		if err != nil {
			return nil, err
		}
		project.TechStack = pq.StringArray(stack)
	}
	if input.IsFeatured != nil {
		project.IsFeatured = *input.IsFeatured
	}
	if input.IsVisible != nil {
		project.IsVisible = *input.IsVisible
	}
	if input.SortOrder != nil {
		project.SortOrder = *input.SortOrder
	}
	if err := validateProject(project); err != nil {
		return nil, err
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
