package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type ProjectHandler struct {
	projectUC usecase.ProjectUseCase
}

func NewProjectHandler(projectUC usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{projectUC: projectUC}
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.projectUC.List(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch projects")
		return
	}

	shared.Success(c, projects)
}

func (h *ProjectHandler) Featured(c *gin.Context) {
	projects, err := h.projectUC.ListFeatured(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to fetch featured projects")
		return
	}

	shared.Success(c, projects)
}
