package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type TagHandler struct {
	tagUC usecase.TagUseCase
}

func NewTagHandler(tagUC usecase.TagUseCase) *TagHandler {
	return &TagHandler{tagUC: tagUC}
}

func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tagUC.List(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch tags")
		return
	}

	shared.Success(c, tags)
}
