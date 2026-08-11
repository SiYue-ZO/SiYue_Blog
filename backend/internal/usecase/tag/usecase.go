package tag

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	tagRepo repo.TagRepo
}

var slugPattern = regexp.MustCompile(`^[\p{L}\p{N}]+(?:-[\p{L}\p{N}]+)*$`)

func NewUseCase(tagRepo repo.TagRepo) usecase.TagUseCase {
	return &useCase{
		tagRepo: tagRepo,
	}
}

func (uc *useCase) List(ctx context.Context) ([]entity.Tag, error) {
	tags, err := uc.tagRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (uc *useCase) Create(ctx context.Context, name, slug string) (*entity.Tag, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))
	if name == "" || !slugPattern.MatchString(slug) {
		return nil, fmt.Errorf("name and a URL-safe slug are required")
	}
	tag := &entity.Tag{
		Name: name,
		Slug: slug,
	}

	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}

	return tag, nil
}

func (uc *useCase) Delete(ctx context.Context, id uint64) error {
	if err := uc.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}
