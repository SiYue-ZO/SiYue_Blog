package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}

	accessToken, refreshToken, user, err := h.authUC.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		shared.Error(c, http.StatusUnauthorized, shared.ErrCodeInvalidCredentials, "Invalid email or password")
		return
	}

	shared.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"avatar_url": user.AvatarURL,
			"role":      user.Role,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}

	user, err := h.authUC.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		shared.Error(c, http.StatusConflict, shared.ErrCodeValidation, err.Error())
		return
	}

	shared.Created(c, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"avatar_url": user.AvatarURL,
		"role":      user.Role,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authUC.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		shared.Error(c, http.StatusUnauthorized, shared.ErrCodeInvalidToken, "Invalid refresh token")
		return
	}

	shared.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
