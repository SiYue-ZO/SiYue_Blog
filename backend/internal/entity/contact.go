package entity

import "time"

type ContactStatus string

const (
	ContactStatusUnread   ContactStatus = "unread"
	ContactStatusRead     ContactStatus = "read"
	ContactStatusArchived ContactStatus = "archived"
)

type ContactMessage struct {
	ID        uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string        `json:"name" gorm:"type:varchar(100);not null"`
	Email     string        `json:"email" gorm:"type:varchar(255);not null"`
	Message   string        `json:"message" gorm:"type:text;not null"`
	Status    ContactStatus `json:"status" gorm:"type:contact_status;not null;default:'unread'"`
	IPHash    string        `json:"-" gorm:"type:varchar(64);not null;default:''"`
	CreatedAt time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ContactMessage) TableName() string { return "contact_messages" }
