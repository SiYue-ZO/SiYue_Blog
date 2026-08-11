package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/pkg/config"
	"siyue-blog/pkg/jwtutil"
)

type authUserRepo struct {
	user    *entity.User
	err     error
	updated bool
}

func (r *authUserRepo) Create(context.Context, *entity.User) error { return nil }
func (r *authUserRepo) GetByID(context.Context, uint64) (*entity.User, error) {
	return r.user, r.err
}
func (r *authUserRepo) GetByEmail(context.Context, string) (*entity.User, error) {
	return r.user, r.err
}
func (r *authUserRepo) GetByUsername(context.Context, string) (*entity.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *authUserRepo) Update(context.Context, *entity.User) error {
	r.updated = true
	return nil
}

type authSessionRepo struct {
	created     *entity.AuthSession
	replacement *entity.AuthSession
	rotateErr   error
	revokeErr   error
	revokedUser uint64
}

func (r *authSessionRepo) Create(_ context.Context, session *entity.AuthSession) error {
	r.created = session
	return nil
}
func (r *authSessionRepo) Rotate(_ context.Context, _ string, replacement *entity.AuthSession) error {
	r.replacement = replacement
	return r.rotateErr
}
func (r *authSessionRepo) Revoke(context.Context, string) error { return r.revokeErr }
func (r *authSessionRepo) RevokeAllForUser(_ context.Context, userID uint64) error {
	r.revokedUser = userID
	return r.revokeErr
}

func authConfig() *config.Config {
	return &config.Config{JWT: config.JWT{
		AccessSecret:  "access-secret-with-at-least-thirty-two-characters",
		RefreshSecret: "refresh-secret-with-at-least-thirty-two-characters",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	}}
}

func activeUser(t *testing.T) (*entity.User, string) {
	t.Helper()
	password := "correct-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return &entity.User{
		ID: 42, Username: "owner", Email: "owner@example.com", PasswordHash: string(hash),
		Role: string(entity.UserRoleAdmin), Status: entity.UserStatusActive,
	}, password
}

func TestLoginCreatesInitialRefreshFamily(t *testing.T) {
	user, password := activeUser(t)
	sessions := &authSessionRepo{}
	uc := NewUseCase(&authUserRepo{user: user}, sessions, authConfig())

	_, refreshToken, _, err := uc.Login(context.Background(), user.Email, password)
	if err != nil {
		t.Fatal(err)
	}
	if sessions.created == nil {
		t.Fatal("refresh session was not created")
	}
	wantHash := hashToken(refreshToken)
	if sessions.created.TokenHash != wantHash || sessions.created.FamilyID != wantHash {
		t.Fatalf("initial session did not create its own family: %#v", sessions.created)
	}
}

func TestRefreshPropagatesSessionReuse(t *testing.T) {
	user, _ := activeUser(t)
	cfg := authConfig()
	refreshToken, err := jwtutil.GenerateRefreshToken(cfg.JWT.RefreshSecret, user.ID, cfg.JWT.RefreshTTL)
	if err != nil {
		t.Fatal(err)
	}
	sessions := &authSessionRepo{rotateErr: repo.ErrSessionReused}
	uc := NewUseCase(&authUserRepo{user: user}, sessions, cfg)

	_, _, err = uc.RefreshToken(context.Background(), refreshToken)
	if !errors.Is(err, repo.ErrSessionReused) {
		t.Fatalf("expected session reuse error, got %v", err)
	}
	if sessions.replacement == nil || sessions.replacement.UserID != user.ID {
		t.Fatal("replacement session was not passed to the repository")
	}
}

func TestLogoutPropagatesRepositoryFailure(t *testing.T) {
	wantErr := errors.New("database unavailable")
	uc := NewUseCase(&authUserRepo{}, &authSessionRepo{revokeErr: wantErr}, authConfig())
	if err := uc.Logout(context.Background(), "refresh-token"); !errors.Is(err, wantErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestChangePasswordRevokesSessionsAndUpdatesHash(t *testing.T) {
	user, currentPassword := activeUser(t)
	users := &authUserRepo{user: user}
	sessions := &authSessionRepo{}
	uc := NewUseCase(users, sessions, authConfig())
	if err := uc.ChangePassword(context.Background(), user.ID, currentPassword, "a-new-secure-password"); err != nil {
		t.Fatal(err)
	}
	if sessions.revokedUser != user.ID || !users.updated {
		t.Fatal("password change did not revoke sessions and persist the user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("a-new-secure-password")); err != nil {
		t.Fatal("new password hash was not stored")
	}
}
