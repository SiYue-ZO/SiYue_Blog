package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/pkg/config"
)

const disabledLegacyAdminPassword = "!disabled-legacy-development-credential!"

// EnsureAdmin creates or synchronizes the explicitly configured bootstrap
// administrator without shipping a known credential in migrations or images.
func EnsureAdmin(ctx context.Context, userRepo repo.UserRepo, cfg config.Admin) error {
	if cfg.Password == "" {
		return nil
	}
	email := strings.ToLower(strings.TrimSpace(cfg.Email))
	username := strings.TrimSpace(cfg.Username)
	if username == "" {
		username = "admin"
	}

	user, err := userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("find bootstrap admin: %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash bootstrap admin password: %w", err)
		}
		if err := userRepo.Create(ctx, &entity.User{
			Username: username, Email: email, PasswordHash: string(hash),
			Role: string(entity.UserRoleAdmin), Status: entity.UserStatusActive,
		}); err != nil {
			return fmt.Errorf("create bootstrap admin: %w", err)
		}
		return nil
	}

	// Existing accounts are never promoted, reactivated, or assigned a new
	// password during normal startup. The only exception is the exact sentinel
	// written by migration 006 for the previously shipped development account.
	if user.PasswordHash != disabledLegacyAdminPassword {
		if user.Role != string(entity.UserRoleAdmin) {
			return fmt.Errorf("bootstrap email belongs to an existing non-admin account")
		}
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash bootstrap admin password: %w", err)
	}
	user.Username = username
	user.PasswordHash = string(hash)
	user.Role = string(entity.UserRoleAdmin)
	user.Status = entity.UserStatusActive
	if err := userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update bootstrap admin: %w", err)
	}
	return nil
}
