package post

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	postRepo repo.PostRepo
	tagRepo  repo.TagRepo
}

var slugPattern = regexp.MustCompile(`^[\p{L}\p{N}]+(?:-[\p{L}\p{N}]+)*$`)

func validateStatus(value entity.PostStatus) bool {
	return value == entity.PostStatusDraft || value == entity.PostStatusPublished || value == entity.PostStatusArchived
}

func validatePost(post *entity.Post) error {
	if post.Title == "" || utf8.RuneCountInString(post.Title) > 255 {
		return fmt.Errorf("title must contain between 1 and 255 characters")
	}
	if !slugPattern.MatchString(post.Slug) || utf8.RuneCountInString(post.Slug) > 255 {
		return fmt.Errorf("a URL-safe slug of at most 255 characters is required")
	}
	if utf8.RuneCountInString(post.Excerpt) > 500 {
		return fmt.Errorf("excerpt must not exceed 500 characters")
	}
	if utf8.RuneCountInString(post.CoverImage) > 512 {
		return fmt.Errorf("cover image URL must not exceed 512 characters")
	}
	if post.CoverImage != "" {
		parsed, err := url.ParseRequestURI(post.CoverImage)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return fmt.Errorf("cover image must be a valid HTTPS URL")
		}
	}
	if !validateStatus(post.Status) {
		return fmt.Errorf("invalid post status")
	}
	if post.Status == entity.PostStatusPublished && strings.TrimSpace(post.Content) == "" {
		return fmt.Errorf("published posts must contain body content")
	}
	return nil
}

func (uc *useCase) resolveTags(ctx context.Context, ids []uint64) ([]entity.Tag, error) {
	if len(ids) > 20 {
		return nil, fmt.Errorf("a post can have at most 20 tags")
	}
	tags := make([]entity.Tag, 0, len(ids))
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, fmt.Errorf("tag IDs must be positive")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		tag, err := uc.tagRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("get tag %d: %w", id, err)
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func NewUseCase(postRepo repo.PostRepo, tagRepo repo.TagRepo) usecase.PostUseCase {
	return &useCase{
		postRepo: postRepo,
		tagRepo:  tagRepo,
	}
}

func (uc *useCase) Create(ctx context.Context, input usecase.CreatePostInput) (*entity.Post, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Slug = strings.ToLower(strings.TrimSpace(input.Slug))
	input.Excerpt = strings.TrimSpace(input.Excerpt)
	input.CoverImage = strings.TrimSpace(input.CoverImage)
	if input.AuthorID == 0 {
		return nil, fmt.Errorf("author is required")
	}
	if input.Status == "" {
		input.Status = string(entity.PostStatusDraft)
	}
	tags, err := uc.resolveTags(ctx, input.TagIDs)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	post := &entity.Post{
		Title:      input.Title,
		Slug:       input.Slug,
		Excerpt:    input.Excerpt,
		Content:    input.Content,
		CoverImage: input.CoverImage,
		Status:     entity.PostStatus(input.Status),
		IsPinned:   input.IsPinned,
		AuthorID:   input.AuthorID,
	}
	if err := validatePost(post); err != nil {
		return nil, err
	}

	if post.Status == entity.PostStatusPublished {
		post.PublishedAt = &now
	}

	if err := uc.postRepo.CreateWithTags(ctx, post, tags); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return post, nil
}

func (uc *useCase) GetByID(ctx context.Context, id uint64) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
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

func (uc *useCase) ListAdmin(ctx context.Context, params usecase.ListAdminPostsParams) ([]entity.Post, int64, error) {
	if params.Status != "" && !validateStatus(entity.PostStatus(params.Status)) {
		return nil, 0, fmt.Errorf("invalid post status")
	}
	posts, total, err := uc.postRepo.ListAdmin(ctx, (params.Page-1)*params.PageSize, params.PageSize, params.Status, strings.TrimSpace(params.Search))
	if err != nil {
		return nil, 0, fmt.Errorf("list admin posts: %w", err)
	}
	return posts, total, nil
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
		post.Title = strings.TrimSpace(*input.Title)
		if post.Title == "" {
			return nil, fmt.Errorf("title is required")
		}
	}
	if input.Slug != nil {
		post.Slug = strings.ToLower(strings.TrimSpace(*input.Slug))
		if !slugPattern.MatchString(post.Slug) {
			return nil, fmt.Errorf("invalid post slug")
		}
	}
	if input.Excerpt != nil {
		post.Excerpt = strings.TrimSpace(*input.Excerpt)
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.CoverImage != nil {
		post.CoverImage = strings.TrimSpace(*input.CoverImage)
	}
	if input.Status != nil {
		newStatus := entity.PostStatus(*input.Status)
		if !validateStatus(newStatus) {
			return nil, fmt.Errorf("invalid post status")
		}
		if newStatus == entity.PostStatusPublished && post.Status != entity.PostStatusPublished {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = newStatus
	}
	if input.IsPinned != nil {
		post.IsPinned = *input.IsPinned
	}
	if err := validatePost(post); err != nil {
		return nil, err
	}

	var tags []entity.Tag
	replaceTags := input.TagIDs != nil
	if input.TagIDs != nil {
		tags, err = uc.resolveTags(ctx, input.TagIDs)
		if err != nil {
			return nil, err
		}
	}

	if err := uc.postRepo.UpdateWithTags(ctx, post, tags, replaceTags); err != nil {
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

func (uc *useCase) IncrementViewCount(ctx context.Context, id uint64) (int, error) {
	return uc.postRepo.IncrementViewCount(ctx, id)
}
