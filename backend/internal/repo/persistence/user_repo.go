package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"siyue-blog/internal/entity"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
