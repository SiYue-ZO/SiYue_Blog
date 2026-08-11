package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
	"siyue-blog/pkg/config"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
	cfg    *config.Config
}

func NewAuthHandler(authUC usecase.AuthUseCase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authUC: authUC, cfg: cfg}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=12,max=128"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,max=128"`
	NewPassword     string `json:"new_password" binding:"required,min=12,max=128"`
}

func (h *AuthHandler) setCookie(c *gin.Context, name, value, path string, maxAge int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(name, value, maxAge, path, h.cfg.Auth.CookieDomain, h.cfg.Auth.CookieSecure, true)
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
	h.setCookie(c, "access_token", accessToken, "/api", int(h.cfg.JWT.AccessTTL.Seconds()))
	h.setCookie(c, "refresh_token", refreshToken, "/api/v1/auth", int(h.cfg.JWT.RefreshTTL.Seconds()))

	shared.Success(c, gin.H{
		"expires_in": int(h.cfg.JWT.AccessTTL.Seconds()),
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
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
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"role":       user.Role,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if refreshToken == "" {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "refresh token is required")
			return
		}
		refreshToken = req.RefreshToken
	}

	accessToken, newRefreshToken, err := h.authUC.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		shared.Error(c, http.StatusUnauthorized, shared.ErrCodeInvalidToken, "Invalid refresh token")
		return
	}
	h.setCookie(c, "access_token", accessToken, "/api", int(h.cfg.JWT.AccessTTL.Seconds()))
	h.setCookie(c, "refresh_token", newRefreshToken, "/api/v1/auth", int(h.cfg.JWT.RefreshTTL.Seconds()))

	shared.Success(c, gin.H{
		"expires_in": int(h.cfg.JWT.AccessTTL.Seconds()),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	shared.Success(c, gin.H{
		"id": c.GetUint64("user_id"), "username": c.GetString("username"), "role": c.GetString("role"),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	if err := h.authUC.Logout(c.Request.Context(), refreshToken); err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "Failed to revoke session")
		return
	}
	h.setCookie(c, "access_token", "", "/api", -1)
	h.setCookie(c, "refresh_token", "", "/api/v1/auth", -1)
	shared.Success(c, gin.H{"logged_out": true})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	if err := h.authUC.ChangePassword(c.Request.Context(), c.GetUint64("user_id"), req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, usecase.ErrInvalidPasswordChange) {
			shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeInternal, "unable to change password")
		return
	}
	h.setCookie(c, "access_token", "", "/api", -1)
	h.setCookie(c, "refresh_token", "", "/api/v1/auth", -1)
	shared.Success(c, gin.H{"changed": true, "login_required": true})
}
