package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
	"siyue-blog/pkg/config"
	"siyue-blog/pkg/jwtutil"
)

type useCase struct {
	userRepo    repo.UserRepo
	sessionRepo repo.SessionRepo
	cfg         *config.Config
}

func NewUseCase(userRepo repo.UserRepo, sessionRepo repo.SessionRepo, cfg *config.Config) usecase.AuthUseCase {
	return &useCase{
		userRepo: userRepo, sessionRepo: sessionRepo, cfg: cfg,
	}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (uc *useCase) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *entity.User, err error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err = uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, fmt.Errorf("invalid credentials")
		}
		return "", "", nil, fmt.Errorf("get user: %w", err)
	}

	if user.Status != entity.UserStatusActive {
		return "", "", nil, fmt.Errorf("account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err = jwtutil.GenerateAccessToken(
		uc.cfg.JWT.AccessSecret,
		user.ID,
		user.Username,
		user.Role,
		uc.cfg.JWT.AccessTTL,
	)
	if err != nil {
		return "", "", nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = jwtutil.GenerateRefreshToken(
		uc.cfg.JWT.RefreshSecret,
		user.ID,
		uc.cfg.JWT.RefreshTTL,
	)
	if err != nil {
		return "", "", nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshHash := hashToken(refreshToken)
	if err := uc.sessionRepo.Create(ctx, &entity.AuthSession{
		UserID: user.ID, TokenHash: refreshHash, FamilyID: refreshHash, ExpiresAt: time.Now().Add(uc.cfg.JWT.RefreshTTL),
	}); err != nil {
		return "", "", nil, fmt.Errorf("create login session: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

func (uc *useCase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	// Check if username already exists
	existingUser, err = uc.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check username: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         string(entity.UserRoleUser),
		Status:       entity.UserStatusActive,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (uc *useCase) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	userID, err := jwtutil.ValidateRefreshToken(uc.cfg.JWT.RefreshSecret, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}

	if user.Status != entity.UserStatusActive {
		return "", "", fmt.Errorf("account is not active")
	}

	newAccessToken, err = jwtutil.GenerateAccessToken(
		uc.cfg.JWT.AccessSecret,
		user.ID,
		user.Username,
		user.Role,
		uc.cfg.JWT.AccessTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err = jwtutil.GenerateRefreshToken(
		uc.cfg.JWT.RefreshSecret,
		user.ID,
		uc.cfg.JWT.RefreshTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	if err := uc.sessionRepo.Rotate(ctx, hashToken(refreshToken), &entity.AuthSession{
		UserID: user.ID, TokenHash: hashToken(newRefreshToken), ExpiresAt: time.Now().Add(uc.cfg.JWT.RefreshTTL),
	}); err != nil {
		return "", "", fmt.Errorf("rotate refresh session: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (uc *useCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return uc.sessionRepo.Revoke(ctx, hashToken(refreshToken))
}

func (uc *useCase) ChangePassword(ctx context.Context, userID uint64, currentPassword, newPassword string) error {
	if userID == 0 {
		return fmt.Errorf("%w: invalid user", usecase.ErrInvalidPasswordChange)
	}
	if len(newPassword) < 12 || len([]byte(newPassword)) > 72 {
		return fmt.Errorf("%w: new password must contain at least 12 characters and at most 72 bytes", usecase.ErrInvalidPasswordChange)
	}
	if currentPassword == newPassword {
		return fmt.Errorf("%w: new password must be different", usecase.ErrInvalidPasswordChange)
	}
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("%w: current password is incorrect", usecase.ErrInvalidPasswordChange)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	// Revoke refresh sessions first. If the following update fails, the old
	// password remains valid and the operation can be retried safely.
	if err := uc.sessionRepo.RevokeAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}
	user.PasswordHash = string(hash)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
