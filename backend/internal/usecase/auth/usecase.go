package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
	"siyue-blog/pkg/config"
	"siyue-blog/pkg/jwtutil"
)

type useCase struct {
	userRepo repo.UserRepo
	cfg      *config.Config
}

func NewUseCase(userRepo repo.UserRepo, cfg *config.Config) usecase.AuthUseCase {
	return &useCase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (uc *useCase) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, user *entity.User, err error) {
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

	return accessToken, refreshToken, user, nil
}

func (uc *useCase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	// Check if email already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	existingUser, _ = uc.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
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

	return newAccessToken, newRefreshToken, nil
}
