package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewConfigReadsEnvironmentOnlySecrets(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
	t.Setenv("SIYUE_JWT_ACCESS_SECRET", "access-secret-with-at-least-thirty-two-characters")
	t.Setenv("SIYUE_JWT_REFRESH_SECRET", "refresh-secret-with-at-least-thirty-two-characters")
	t.Setenv("SIYUE_ADMIN_EMAIL", "admin@example.com")
	t.Setenv("SIYUE_ADMIN_PASSWORD", "a-strong-password")
	cfg, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.JWT.AccessSecret == "" || cfg.Admin.Email != "admin@example.com" {
		t.Fatalf("environment values were not loaded: %#v", cfg)
	}
}

func validConfig() Config {
	return Config{
		App:  App{Env: "production"},
		HTTP: HTTP{ReadTimeout: time.Second, WriteTimeout: time.Second, IdleTimeout: time.Second, MaxBodyBytes: 1024},
		Postgres: Postgres{
			Host: "postgres", Port: 5432, User: "app", Password: "database-password-38f5c94d", DBName: "blog", SSLMode: "disable",
			MaxOpenConns: 25, MaxIdleConns: 10, MaxLifetime: time.Minute,
		},
		JWT: JWT{
			AccessSecret:  "access-secret-with-at-least-thirty-two-characters",
			RefreshSecret: "refresh-secret-with-at-least-thirty-two-characters",
			AccessTTL:     15 * time.Minute, RefreshTTL: 24 * time.Hour,
		},
		Auth:  Auth{CookieSecure: true},
		Admin: Admin{Username: "admin", Email: "admin@siyue.dev", Password: "strong-admin-password-38f5c94d"},
		CORS:  CORS{AllowOrigins: []string{"https://blog.example.com"}},
	}
}

func TestValidateRejectsProductionPlaceholders(t *testing.T) {
	cfg := validConfig()
	cfg.JWT.AccessSecret = "replace-with-at-least-thirty-two-random-characters"
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "placeholder") {
		t.Fatalf("expected placeholder validation error, got %v", err)
	}
}

func TestValidateRejectsOversizedBcryptPassword(t *testing.T) {
	cfg := validConfig()
	cfg.Admin.Password = strings.Repeat("密", 25)
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "72 bytes") {
		t.Fatalf("expected bcrypt length validation error, got %v", err)
	}
}

func TestValidateProductionConfig(t *testing.T) {
	cfg := validConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid configuration was rejected: %v", err)
	}
}

func TestValidateRejectsSharedSecrets(t *testing.T) {
	cfg := validConfig()
	cfg.JWT.RefreshSecret = cfg.JWT.AccessSecret
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "different") {
		t.Fatalf("expected shared-secret validation error, got %v", err)
	}
}

func TestValidateRejectsInsecureProductionCookie(t *testing.T) {
	cfg := validConfig()
	cfg.Auth.CookieSecure = false
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "cookie_secure") {
		t.Fatalf("expected cookie validation error, got %v", err)
	}
}
