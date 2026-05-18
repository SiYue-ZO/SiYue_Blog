package entity

import (
	"time"

	"github.com/lib/pq"
)

type Project struct {
	ID          uint64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string      `json:"title" gorm:"type:varchar(255);not null"`
	Description string      `json:"description" gorm:"type:text"`
	RepoURL     string      `json:"repo_url" gorm:"type:varchar(512)"`
	DemoURL     string      `json:"demo_url" gorm:"type:varchar(512)"`
	ImageURL    string      `json:"image_url" gorm:"type:varchar(512)"`
	TechStack   pq.StringArray `json:"tech_stack" gorm:"type:text[]"`
	IsFeatured  bool        `json:"is_featured" gorm:"not null;default:false"`
	SortOrder   int         `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Project) TableName() string {
	return "projects"
}
