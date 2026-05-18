package entity

import "time"

type User struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string     `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
	Email        string     `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null"`
	AvatarURL    string     `json:"avatar_url" gorm:"type:varchar(512)"`
	Role         string     `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	Status       UserStatus `json:"status" gorm:"type:user_status;not null;default:'active'"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
