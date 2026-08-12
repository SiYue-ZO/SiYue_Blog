package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

type HomePageHandler struct {
	homePageUC usecase.HomePageUseCase
}

func NewHomePageHandler(homePageUC usecase.HomePageUseCase) *HomePageHandler {
	return &HomePageHandler{homePageUC: homePageUC}
}

func (h *HomePageHandler) GetPublished(c *gin.Context) {
	config, err := h.homePageUC.GetPublished(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to load home page")
		return
	}
	shared.Success(c, config)
}

func (h *HomePageHandler) GetAdmin(c *gin.Context) {
	content, err := h.homePageUC.GetAdmin(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to load home page")
		return
	}
	shared.Success(c, content)
}

func (h *HomePageHandler) SaveDraft(c *gin.Context) {
	var config entity.HomePageConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid home page configuration")
		return
	}
	content, err := h.homePageUC.SaveDraft(c.Request.Context(), config, c.GetUint64("user_id"))
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, content)
}

func (h *HomePageHandler) Publish(c *gin.Context) {
	content, err := h.homePageUC.Publish(c.Request.Context(), c.GetUint64("user_id"))
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, content)
}

func (h *HomePageHandler) ResetDraft(c *gin.Context) {
	content, err := h.homePageUC.ResetDraft(c.Request.Context(), c.GetUint64("user_id"))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to reset home page draft")
		return
	}
	shared.Success(c, content)
}
