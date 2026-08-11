package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/entity"
	"siyue-blog/pkg/config"
)

type logoutAuthUseCase struct{ logoutErr error }

func (f *logoutAuthUseCase) Login(context.Context, string, string) (string, string, *entity.User, error) {
	return "", "", nil, nil
}
func (f *logoutAuthUseCase) Register(context.Context, string, string, string) (*entity.User, error) {
	return nil, nil
}
func (f *logoutAuthUseCase) RefreshToken(context.Context, string) (string, string, error) {
	return "", "", nil
}
func (f *logoutAuthUseCase) Logout(context.Context, string) error { return f.logoutErr }
func (f *logoutAuthUseCase) ChangePassword(context.Context, uint64, string, string) error {
	return nil
}

func TestLogoutReturnsServerErrorWhenSessionCannotBeRevoked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&logoutAuthUseCase{logoutErr: errors.New("database unavailable")}, &config.Config{})
	router := gin.New()
	router.POST("/logout", handler.Logout)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/logout", nil)
	request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if cookies := recorder.Header().Values("Set-Cookie"); len(cookies) != 0 {
		t.Fatalf("cookies were cleared even though revocation failed: %v", cookies)
	}
}
