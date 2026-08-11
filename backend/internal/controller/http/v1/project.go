package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type ProjectHandler struct {
	projectUC usecase.ProjectUseCase
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var input usecase.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	project, err := h.projectUC.Create(c.Request.Context(), input)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Created(c, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid project ID")
		return
	}
	var input usecase.UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	project, err := h.projectUC.Update(c.Request.Context(), id, input)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid project ID")
		return
	}
	if err := h.projectUC.Delete(c.Request.Context(), id); err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodeProjectNotFound, "project not found")
		return
	}
	shared.Success(c, gin.H{"deleted": true})
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

func (h *ProjectHandler) AdminList(c *gin.Context) {
	projects, err := h.projectUC.ListAdmin(c.Request.Context())
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
