package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

type SettingsHandler struct {
	settingsUC usecase.SettingsUseCase
}

func NewSettingsHandler(settingsUC usecase.SettingsUseCase) *SettingsHandler {
	return &SettingsHandler{settingsUC: settingsUC}
}

func (h *SettingsHandler) Get(c *gin.Context) {
	settings, err := h.settingsUC.Get(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to load site settings")
		return
	}
	shared.Success(c, settings)
}

func (h *SettingsHandler) Update(c *gin.Context) {
	var input struct {
		ProjectPageVisible *bool                `json:"project_page_visible"`
		SocialLinks        *[]entity.SocialLink `json:"social_links"`
		Footer             *entity.FooterConfig `json:"footer"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || (input.ProjectPageVisible == nil && input.SocialLinks == nil && input.Footer == nil) {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "at least one setting is required")
		return
	}
	settings, err := h.settingsUC.Update(c.Request.Context(), input.ProjectPageVisible, input.SocialLinks, input.Footer)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, settings)
}
