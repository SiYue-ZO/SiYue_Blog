package bootstrap

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/pkg/config"
)

type adminUserRepo struct {
	user      *entity.User
	getErr    error
	created   *entity.User
	updated   *entity.User
	createErr error
	updateErr error
}

func (r *adminUserRepo) Create(_ context.Context, user *entity.User) error {
	r.created = user
	return r.createErr
}

func (r *adminUserRepo) GetByID(context.Context, uint64) (*entity.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *adminUserRepo) GetByEmail(context.Context, string) (*entity.User, error) {
	return r.user, r.getErr
}

func (r *adminUserRepo) GetByUsername(context.Context, string) (*entity.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *adminUserRepo) Update(_ context.Context, user *entity.User) error {
	r.updated = user
	return r.updateErr
}

func bootstrapConfig() config.Admin {
	return config.Admin{Username: "owner", Email: "owner@example.com", Password: "new-bootstrap-password"}
}

func TestEnsureAdminCreatesMissingAdministrator(t *testing.T) {
	repository := &adminUserRepo{getErr: gorm.ErrRecordNotFound}
	if err := EnsureAdmin(context.Background(), repository, bootstrapConfig()); err != nil {
		t.Fatal(err)
	}
	if repository.created == nil {
		t.Fatal("administrator was not created")
	}
	if repository.created.Role != string(entity.UserRoleAdmin) || repository.created.Status != entity.UserStatusActive {
		t.Fatalf("unexpected administrator: %#v", repository.created)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.created.PasswordHash), []byte(bootstrapConfig().Password)); err != nil {
		t.Fatalf("created password does not match configured password: %v", err)
	}
}

func TestEnsureAdminLeavesExistingAdministratorUnchanged(t *testing.T) {
	original := &entity.User{
		ID: 7, Username: "existing", Email: "owner@example.com", PasswordHash: "existing-hash",
		Role: string(entity.UserRoleAdmin), Status: entity.UserStatusBanned,
	}
	repository := &adminUserRepo{user: original}
	if err := EnsureAdmin(context.Background(), repository, bootstrapConfig()); err != nil {
		t.Fatal(err)
	}
	if repository.updated != nil {
		t.Fatal("existing administrator was unexpectedly updated")
	}
	if original.PasswordHash != "existing-hash" || original.Status != entity.UserStatusBanned || original.Username != "existing" {
		t.Fatalf("existing administrator was mutated: %#v", original)
	}
}

func TestEnsureAdminTakesOverDisabledLegacySentinel(t *testing.T) {
	legacy := &entity.User{
		ID: 1, Username: "admin", Email: "owner@example.com", PasswordHash: disabledLegacyAdminPassword,
		Role: string(entity.UserRoleAdmin), Status: entity.UserStatusInactive,
	}
	repository := &adminUserRepo{user: legacy}
	if err := EnsureAdmin(context.Background(), repository, bootstrapConfig()); err != nil {
		t.Fatal(err)
	}
	if repository.updated == nil {
		t.Fatal("legacy administrator was not taken over")
	}
	if repository.updated.Username != "owner" || repository.updated.Role != string(entity.UserRoleAdmin) || repository.updated.Status != entity.UserStatusActive {
		t.Fatalf("unexpected migrated administrator: %#v", repository.updated)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.updated.PasswordHash), []byte(bootstrapConfig().Password)); err != nil {
		t.Fatalf("migrated password does not match configured password: %v", err)
	}
}

func TestEnsureAdminRefusesExistingNonAdminAccount(t *testing.T) {
	existing := &entity.User{
		ID: 9, Username: "member", Email: "owner@example.com", PasswordHash: "member-hash",
		Role: string(entity.UserRoleUser), Status: entity.UserStatusActive,
	}
	repository := &adminUserRepo{user: existing}
	err := EnsureAdmin(context.Background(), repository, bootstrapConfig())
	if err == nil {
		t.Fatal("expected an existing-account conflict")
	}
	if repository.updated != nil || repository.created != nil {
		t.Fatal("existing non-admin account was modified")
	}
	if existing.PasswordHash != "member-hash" || existing.Role != string(entity.UserRoleUser) {
		t.Fatalf("existing non-admin account was mutated: %#v", existing)
	}
}
