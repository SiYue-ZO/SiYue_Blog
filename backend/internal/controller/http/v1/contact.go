package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/usecase"
	contactusecase "siyue-blog/internal/usecase/contact"
)

type ContactHandler struct {
	contactUC usecase.ContactUseCase
	hashSalt  string
}

func NewContactHandler(contactUC usecase.ContactUseCase, hashSalt string) *ContactHandler {
	return &ContactHandler{contactUC: contactUC, hashSalt: hashSalt}
}

type contactRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Email   string `json:"email" binding:"required,email,max=255"`
	Message string `json:"message" binding:"required,min=10,max=5000"`
	Website string `json:"website" binding:"max=200"`
}

func (h *ContactHandler) Create(c *gin.Context) {
	var req contactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	// Honeypot field: bots receive a normal response without storing spam.
	if req.Website != "" {
		shared.Created(c, gin.H{"received": true})
		return
	}
	digest := sha256.Sum256([]byte(h.hashSalt + c.ClientIP()))
	_, err := h.contactUC.Create(c.Request.Context(), usecase.CreateContactInput{
		Name: req.Name, Email: req.Email, Message: req.Message, IPHash: hex.EncodeToString(digest[:]),
	})
	if err != nil {
		if errors.Is(err, contactusecase.ErrInvalidInput) {
			shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeInternal, "unable to save message")
		return
	}
	shared.Created(c, gin.H{"received": true})
}

func (h *ContactHandler) List(c *gin.Context) {
	params := shared.ParsePagination(c)
	messages, total, err := h.contactUC.List(c.Request.Context(), params.Page, params.PageSize, c.Query("status"))
	if err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Paginated(c, messages, params.Page, params.PageSize, int(total))
}

func (h *ContactHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid message ID")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	if err := h.contactUC.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, err.Error())
		return
	}
	shared.Success(c, gin.H{"updated": true})
}

func (h *ContactHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		shared.Error(c, http.StatusBadRequest, shared.ErrCodeValidation, "invalid message ID")
		return
	}
	if err := h.contactUC.Delete(c.Request.Context(), id); err != nil {
		shared.Error(c, http.StatusNotFound, shared.ErrCodeNotFound, "message not found")
		return
	}
	shared.Success(c, gin.H{"deleted": true})
}
