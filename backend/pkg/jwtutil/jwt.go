package jwtutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func tokenID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate token ID: %w", err)
	}
	return hex.EncodeToString(value), nil
}

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	TokenUse string `json:"token_use"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(secret string, userID uint64, username, role string, ttl time.Duration) (string, error) {
	id, err := tokenID()
	if err != nil {
		return "", err
	}
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "siyue-blog",
			Audience:  jwt.ClaimStrings{"siyue-blog"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(secret string, userID uint64, ttl time.Duration) (string, error) {
	id, err := tokenID()
	if err != nil {
		return "", err
	}
	claims := Claims{
		UserID:   userID,
		TokenUse: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "siyue-blog",
			Audience:  jwt.ClaimStrings{"siyue-blog"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer("siyue-blog"), jwt.WithAudience("siyue-blog"))
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenUse != "access" {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidateRefreshToken(secret, tokenStr string) (uint64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer("siyue-blog"), jwt.WithAudience("siyue-blog"))
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.TokenUse != "refresh" {
		return 0, fmt.Errorf("invalid token")
	}

	if claims.UserID == 0 {
		return 0, fmt.Errorf("invalid token subject")
	}
	return claims.UserID, nil
}
