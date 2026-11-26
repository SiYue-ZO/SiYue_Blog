package post

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

// MockPostRepo is a mock implementation of PostRepo
type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) CreateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag) error {
	args := m.Called(ctx, post, tags)
	return args.Error(0)
}

func (m *MockPostRepo) GetByID(ctx context.Context, id uint64) (*entity.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepo) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostRepo) List(ctx context.Context, offset, limit int, tag string, search string) ([]entity.Post, int64, error) {
	args := m.Called(ctx, offset, limit, tag, search)
	return args.Get(0).([]entity.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepo) ListAdmin(ctx context.Context, offset, limit int, status string, search string) ([]entity.Post, int64, error) {
	args := m.Called(ctx, offset, limit, status, search)
	return args.Get(0).([]entity.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepo) ListFeatured(ctx context.Context, limit int) ([]entity.Post, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, post *entity.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepo) UpdateWithTags(ctx context.Context, post *entity.Post, tags []entity.Tag, replaceTags bool) error {
	args := m.Called(ctx, post, tags, replaceTags)
	return args.Error(0)
}

func (m *MockPostRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepo) RecordView(ctx context.Context, id uint64, eventID string) (int, error) {
	args := m.Called(ctx, id, eventID)
	return args.Int(0), args.Error(1)
}

// MockTagRepo is a mock implementation of TagRepo
type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Create(ctx context.Context, tag *entity.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepo) GetByID(ctx context.Context, id uint64) (*entity.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}

func (m *MockTagRepo) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}

func (m *MockTagRepo) List(ctx context.Context) ([]entity.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Tag), args.Error(1)
}

func (m *MockTagRepo) ListByPostID(ctx context.Context, postID uint64) ([]entity.Tag, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]entity.Tag), args.Error(1)
}

func (m *MockTagRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreatePost(t *testing.T) {
	mockPostRepo := new(MockPostRepo)
	mockTagRepo := new(MockTagRepo)
	uc := NewUseCase(mockPostRepo, mockTagRepo)

	ctx := context.Background()
	now := time.Now()

	input := usecase.CreatePostInput{
		Title:    "Test Post",
		Slug:     "test-post",
		Excerpt:  "This is a test post",
		Content:  "Test content",
		Status:   "published",
		AuthorID: 1,
	}

	mockPostRepo.On("CreateWithTags", ctx, mock.AnythingOfType("*entity.Post"), []entity.Tag{}).
		Run(func(args mock.Arguments) {
			post := args.Get(1).(*entity.Post)
			post.ID = 1
			post.CreatedAt = now
			post.UpdatedAt = now
		}).
		Return(nil)

	result, err := uc.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Post", result.Title)
	assert.Equal(t, "test-post", result.Slug)
	assert.Equal(t, entity.PostStatusPublished, result.Status)
	assert.NotNil(t, result.PublishedAt)
	mockPostRepo.AssertExpectations(t)
}

func TestCreateRejectsPublishedPostWithoutContent(t *testing.T) {
	uc := NewUseCase(new(MockPostRepo), new(MockTagRepo))
	result, err := uc.Create(context.Background(), usecase.CreatePostInput{
		Title: "Empty", Slug: "empty", Status: "published", AuthorID: 1,
	})
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "must contain body content")
}

func TestGetPostBySlug(t *testing.T) {
	mockPostRepo := new(MockPostRepo)
	mockTagRepo := new(MockTagRepo)
	uc := NewUseCase(mockPostRepo, mockTagRepo)

	ctx := context.Background()
	expectedPost := &entity.Post{
		ID:    1,
		Title: "Test Post",
		Slug:  "test-post",
	}

	mockPostRepo.On("GetBySlug", ctx, "test-post").Return(expectedPost, nil)

	result, err := uc.GetBySlug(ctx, "test-post")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPost.Title, result.Title)
	mockPostRepo.AssertExpectations(t)
}

func TestListPosts(t *testing.T) {
	mockPostRepo := new(MockPostRepo)
	mockTagRepo := new(MockTagRepo)
	uc := NewUseCase(mockPostRepo, mockTagRepo)

	ctx := context.Background()
	expectedPosts := []entity.Post{
		{ID: 1, Title: "Post 1", Slug: "post-1"},
		{ID: 2, Title: "Post 2", Slug: "post-2"},
	}

	mockPostRepo.On("List", ctx, 0, 10, "", "").Return(expectedPosts, int64(2), nil)

	result, total, err := uc.List(ctx, usecase.ListPostsParams{
		Page:     1,
		PageSize: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	mockPostRepo.AssertExpectations(t)
}

func TestListFeaturedPosts(t *testing.T) {
	mockPostRepo := new(MockPostRepo)
	mockTagRepo := new(MockTagRepo)
	uc := NewUseCase(mockPostRepo, mockTagRepo)

	ctx := context.Background()
	expectedPosts := []entity.Post{
		{ID: 1, Title: "Featured Post", Slug: "featured-post", IsPinned: true},
	}

	mockPostRepo.On("ListFeatured", ctx, 3).Return(expectedPosts, nil)

	result, err := uc.ListFeatured(ctx, 3)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].IsPinned)
	mockPostRepo.AssertExpectations(t)
}

func TestDeletePost(t *testing.T) {
	mockPostRepo := new(MockPostRepo)
	mockTagRepo := new(MockTagRepo)
	uc := NewUseCase(mockPostRepo, mockTagRepo)

	ctx := context.Background()

	mockPostRepo.On("Delete", ctx, uint64(1)).Return(nil)

	err := uc.Delete(ctx, 1)

	assert.NoError(t, err)
	mockPostRepo.AssertExpectations(t)
}
