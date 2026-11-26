package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

type PostHandler struct {
	postUC usecase.PostUseCase
	log    zerolog.Logger
}

func NewPostHandler(postUC usecase.PostUseCase, log zerolog.Logger) *PostHandler {
	return &PostHandler{postUC: postUC, log: log}
}

type publicAuthor struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type publicPost struct {
	ID          uint64        `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Excerpt     string        `json:"excerpt"`
	Content     string        `json:"content,omitempty"`
	CoverImage  string        `json:"cover_image"`
	ViewCount   int           `json:"view_count"`
	IsPinned    bool          `json:"is_pinned"`
	Author      *publicAuthor `json:"author,omitempty"`
	Tags        []entity.Tag  `json:"tags"`
	PublishedAt *time.Time    `json:"published_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func toPublicPost(post entity.Post, includeContent bool) publicPost {
	result := publicPost{
		ID: post.ID, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt,
		CoverImage: post.CoverImage, ViewCount: post.ViewCount, IsPinned: post.IsPinned,
		Tags: post.Tags, PublishedAt: post.PublishedAt, UpdatedAt: post.UpdatedAt,
	}
	if includeContent {
		result.Content = post.Content
	}
	if post.Author != nil {
		result.Author = &publicAuthor{Username: post.Author.Username, AvatarURL: post.Author.AvatarURL}
	}
	return result
}

func toPublicPosts(posts []entity.Post) []publicPost {
	result := make([]publicPost, len(posts))
	for i := range posts {
		result[i] = toPublicPost(posts[i], false)
	}
	return result
}

func (h *PostHandler) List(c *gin.Context) {
	params := shared.ParsePagination(c)
	tag := c.Query("tag")
	search := c.Query("search")

	posts, total, err := h.postUC.List(c.Request.Context(), usecase.ListPostsParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Tag:      tag,
		Search:   search,
	})
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch posts")
		return
	}

	shared.Paginated(c, toPublicPosts(posts), params.Page, params.PageSize, int(total))
}

func (h *PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeBadRequest, "Slug is required")
		return
	}

	post, err := h.postUC.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodePostNotFound, "Post not found")
		return
	}

	shared.Success(c, toPublicPost(*post, true))
}

// RecordView is separate from GET so metadata, prefetches and server rendering
// cannot count as reading an article. Repeated event IDs are deduplicated in SQL.
func (h *PostHandler) RecordView(c *gin.Context) {
	var input struct {
		EventID string `json:"event_id" binding:"required,uuid4"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "A UUID v4 event_id is required")
		return
	}
	post, err := h.postUC.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodePostNotFound, "Post not found")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	viewCount, err := h.postUC.RecordView(ctx, post.ID, input.EventID)
	if err != nil {
		h.log.Error().Err(err).Uint64("post_id", post.ID).Msg("failed to record post view")
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to record post view")
		return
	}
	c.Header("Cache-Control", "no-store")
	shared.Success(c, gin.H{"view_count": viewCount})
}

func (h *PostHandler) Featured(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 12 {
		limit = 3
	}

	posts, err := h.postUC.ListFeatured(c.Request.Context(), limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch featured posts")
		return
	}

	shared.Success(c, toPublicPosts(posts))
}

func (h *PostHandler) AdminList(c *gin.Context) {
	params := shared.ParsePagination(c)
	posts, total, err := h.postUC.ListAdmin(c.Request.Context(), usecase.ListAdminPostsParams{
		Page: params.Page, PageSize: params.PageSize, Status: c.Query("status"), Search: c.Query("search"),
	})
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Paginated(c, posts, params.Page, params.PageSize, int(total))
}

func (h *PostHandler) AdminGet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid post ID")
		return
	}
	post, err := h.postUC.GetByID(c.Request.Context(), id)
	if err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodePostNotFound, "post not found")
		return
	}
	shared.Success(c, post)
}

func (h *PostHandler) Create(c *gin.Context) {
	var input usecase.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	input.AuthorID = c.GetUint64("user_id")
	post, err := h.postUC.Create(c.Request.Context(), input)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Created(c, post)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid post ID")
		return
	}
	var input usecase.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	post, err := h.postUC.Update(c.Request.Context(), id, input)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid post ID")
		return
	}
	if err := h.postUC.Delete(c.Request.Context(), id); err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodePostNotFound, "post not found")
		return
	}
	shared.Success(c, gin.H{"deleted": true})
}
