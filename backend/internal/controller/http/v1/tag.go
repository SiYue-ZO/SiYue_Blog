package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type TagHandler struct {
	tagUC usecase.TagUseCase
}

type createTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=50"`
}

func (h *TagHandler) Create(c *gin.Context) {
	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	tag, err := h.tagUC.Create(c.Request.Context(), strings.TrimSpace(req.Name), strings.ToLower(strings.TrimSpace(req.Slug)))
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Created(c, tag)
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid tag ID")
		return
	}
	if err := h.tagUC.Delete(c.Request.Context(), id); err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodeTagNotFound, "tag not found")
		return
	}
	shared.Success(c, gin.H{"deleted": true})
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
