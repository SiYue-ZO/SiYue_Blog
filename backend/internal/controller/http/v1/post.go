package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type PostHandler struct {
	postUC usecase.PostUseCase
}

func NewPostHandler(postUC usecase.PostUseCase) *PostHandler {
	return &PostHandler{postUC: postUC}
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

	shared.Paginated(c, posts, params.Page, params.PageSize, int(total))
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

	// Increment view count
	go h.postUC.IncrementViewCount(c.Request.Context(), post.ID)

	shared.Success(c, post)
}

func (h *PostHandler) Featured(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 10 {
		limit = 3
	}

	posts, err := h.postUC.ListFeatured(c.Request.Context(), limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch featured posts")
		return
	}

	shared.Success(c, posts)
}
