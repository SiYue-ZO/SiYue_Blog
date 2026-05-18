package entity

import "time"

type Tag struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	Slug      string    `json:"slug" gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relations
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}
