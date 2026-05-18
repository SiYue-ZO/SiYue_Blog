package entity

import "time"

type Post struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Slug        string     `json:"slug" gorm:"type:varchar(255);not null;uniqueIndex"`
	Excerpt     string     `json:"excerpt" gorm:"type:varchar(500)"`
	Content     string     `json:"content,omitempty" gorm:"type:text"`
	CoverImage  string     `json:"cover_image" gorm:"type:varchar(512)"`
	Status      PostStatus `json:"status" gorm:"type:post_status;not null;default:'draft'"`
	ViewCount   int        `json:"view_count" gorm:"not null;default:0"`
	IsPinned    bool       `json:"is_pinned" gorm:"not null;default:false"`
	AuthorID    uint64     `json:"author_id" gorm:"not null;index"`
	PublishedAt *time.Time `json:"published_at" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Author *User `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Tags   []Tag `json:"tags,omitempty" gorm:"many2many:post_tags;"`
}

func (Post) TableName() string {
	return "posts"
}
