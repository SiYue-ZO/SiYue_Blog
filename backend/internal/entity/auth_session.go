package entity

import "time"

type AuthSession struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint64     `json:"user_id" gorm:"not null;index"`
	TokenHash  string     `json:"-" gorm:"type:varchar(64);not null;uniqueIndex"`
	FamilyID   string     `json:"-" gorm:"type:varchar(64);not null;index"`
	ReplacedBy string     `json:"-" gorm:"type:varchar(64)"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (AuthSession) TableName() string { return "auth_sessions" }
