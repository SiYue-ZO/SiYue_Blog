package jwtutil

import (
	"testing"
	"time"
)

func TestAccessAndRefreshTokensAreNotInterchangeable(t *testing.T) {
	accessSecret := "access-secret-with-at-least-thirty-two-characters"
	refreshSecret := "refresh-secret-with-at-least-thirty-two-characters"
	access, err := GenerateAccessToken(accessSecret, 42, "admin", "admin", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateToken(accessSecret, access)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 42 || claims.TokenUse != "access" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
	if _, err := ValidateRefreshToken(accessSecret, access); err == nil {
		t.Fatal("access token was accepted as a refresh token")
	}

	refresh, err := GenerateRefreshToken(refreshSecret, 42, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	userID, err := ValidateRefreshToken(refreshSecret, refresh)
	if err != nil {
		t.Fatal(err)
	}
	if userID != 42 {
		t.Fatalf("unexpected user ID: %d", userID)
	}
	if _, err := ValidateToken(refreshSecret, refresh); err == nil {
		t.Fatal("refresh token was accepted as an access token")
	}
}

func TestExpiredTokenIsRejected(t *testing.T) {
	secret := "access-secret-with-at-least-thirty-two-characters"
	token, err := GenerateAccessToken(secret, 1, "admin", "admin", -time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateToken(secret, token); err == nil {
		t.Fatal("expired token was accepted")
	}
}
