package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"siyue-blog/internal/controller/http/shared"
	"siyue-blog/internal/entity"
)

type AdminHandler struct{ db *gorm.DB }

func NewAdminHandler(db *gorm.DB) *AdminHandler { return &AdminHandler{db: db} }

func (h *AdminHandler) Stats(c *gin.Context) {
	type stats struct {
		TotalPosts     int64 `json:"total_posts"`
		PublishedPosts int64 `json:"published_posts"`
		DraftPosts     int64 `json:"draft_posts"`
		TotalTags      int64 `json:"total_tags"`
		TotalProjects  int64 `json:"total_projects"`
		TotalViews     int64 `json:"total_views"`
		UnreadMessages int64 `json:"unread_messages"`
	}
	var result stats
	queries := []struct {
		dest  *int64
		query *gorm.DB
	}{
		{&result.TotalPosts, h.db.WithContext(c.Request.Context()).Model(&entity.Post{})},
		{&result.PublishedPosts, h.db.WithContext(c.Request.Context()).Model(&entity.Post{}).Where("status = ?", entity.PostStatusPublished)},
		{&result.DraftPosts, h.db.WithContext(c.Request.Context()).Model(&entity.Post{}).Where("status = ?", entity.PostStatusDraft)},
		{&result.TotalTags, h.db.WithContext(c.Request.Context()).Model(&entity.Tag{})},
		{&result.TotalProjects, h.db.WithContext(c.Request.Context()).Model(&entity.Project{})},
		{&result.UnreadMessages, h.db.WithContext(c.Request.Context()).Model(&entity.ContactMessage{}).Where("status = ?", entity.ContactStatusUnread)},
	}
	for _, item := range queries {
		if err := item.query.Count(item.dest).Error; err != nil {
			shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to load dashboard statistics")
			return
		}
	}
	if err := h.db.WithContext(c.Request.Context()).Model(&entity.Post{}).Select("COALESCE(SUM(view_count), 0)").Scan(&result.TotalViews).Error; err != nil {
		shared.Error(c, http.StatusInternalServerError, shared.ErrCodeDatabase, "failed to load dashboard statistics")
		return
	}
	shared.Success(c, result)
}
