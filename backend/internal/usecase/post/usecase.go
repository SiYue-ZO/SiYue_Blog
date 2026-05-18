package post

import (
	"context"
	"fmt"
	"time"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	postRepo repo.PostRepo
	tagRepo  repo.TagRepo
}

func NewUseCase(postRepo repo.PostRepo, tagRepo repo.TagRepo) usecase.PostUseCase {
	return &useCase{
		postRepo: postRepo,
		tagRepo:  tagRepo,
	}
}

func (uc *useCase) Create(ctx context.Context, input usecase.CreatePostInput) (*entity.Post, error) {
	now := time.Now()
	post := &entity.Post{
		Title:      input.Title,
		Slug:       input.Slug,
		Excerpt:    input.Excerpt,
		Content:    input.Content,
		CoverImage: input.CoverImage,
		Status:     entity.PostStatus(input.Status),
		IsPinned:   input.IsPinned,
		AuthorID:   1, // TODO: get from context
	}

	if post.Status == entity.PostStatusPublished {
		post.PublishedAt = &now
	}

	if err := uc.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	// Attach tags
	if len(input.TagIDs) > 0 {
		var tags []entity.Tag
		for _, tagID := range input.TagIDs {
			tag, err := uc.tagRepo.GetByID(ctx, tagID)
			if err != nil {
				continue
			}
			tags = append(tags, *tag)
		}
		post.Tags = tags
		if err := uc.postRepo.Update(ctx, post); err != nil {
			return nil, fmt.Errorf("attach tags: %w", err)
		}
	}

	return post, nil
}

func (uc *useCase) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	post, err := uc.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get post by slug: %w", err)
	}
	return post, nil
}

func (uc *useCase) List(ctx context.Context, params usecase.ListPostsParams) ([]entity.Post, int64, error) {
	offset := (params.Page - 1) * params.PageSize
	posts, total, err := uc.postRepo.List(ctx, offset, params.PageSize, params.Tag, params.Search)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	return posts, total, nil
}

func (uc *useCase) ListFeatured(ctx context.Context, limit int) ([]entity.Post, error) {
	posts, err := uc.postRepo.ListFeatured(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list featured posts: %w", err)
	}
	return posts, nil
}

func (uc *useCase) Update(ctx context.Context, id uint64, input usecase.UpdatePostInput) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Slug != nil {
		post.Slug = *input.Slug
	}
	if input.Excerpt != nil {
		post.Excerpt = *input.Excerpt
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.CoverImage != nil {
		post.CoverImage = *input.CoverImage
	}
	if input.Status != nil {
		newStatus := entity.PostStatus(*input.Status)
		if newStatus == entity.PostStatusPublished && post.Status != entity.PostStatusPublished {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = newStatus
	}
	if input.IsPinned != nil {
		post.IsPinned = *input.IsPinned
	}

	// Update tags if provided
	if input.TagIDs != nil {
		var tags []entity.Tag
		for _, tagID := range input.TagIDs {
			tag, err := uc.tagRepo.GetByID(ctx, tagID)
			if err != nil {
				continue
			}
			tags = append(tags, *tag)
		}
		post.Tags = tags
	}

	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}

	return post, nil
}

func (uc *useCase) Delete(ctx context.Context, id uint64) error {
	if err := uc.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}

func (uc *useCase) IncrementViewCount(ctx context.Context, id uint64) error {
	return uc.postRepo.IncrementViewCount(ctx, id)
}
